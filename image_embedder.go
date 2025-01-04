package stegano

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
)

// EncodeAndSave embeds the provided data into the given image and saves the modified image to a new file.
// The data is embedded using the specified bit depth. If `defaultCompression` is true, the data is compressed before embedding.
// Returns an error if the data exceeds the embedding capacity of the image or if the saving process fails.

// Parameters:
// - coverImage: The original image where data will be embedded.
// - data: The data to embed into the image.
// - bitDepth: The number of bits per channel used for embedding (0-7).
// - outputFilename: The name of the file where the modified image will be saved.
// - defaultCompression: A flag indicating whether the data should be compressed before embedding.
func (m *EmbedHandler) Encode(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	// Validate coverImage dimensions
	if coverImage == nil {
		return ErrInvalidCoverImage
	}
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()
	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	// Validate bit depth
	if bitDepth < 0 || bitDepth > 7 {
		return ErrDepthOutOfRange
	}

	// Validate data
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}
	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return ErrFailedToExtractRGB
	}

	maxCapacity := (len(RGBchannels) * 3 * (int(bitDepth) + 1)) / 8
	if (len(data)*8)+32 > maxCapacity {
		return fmt.Errorf("data is too large to embed into the image: maxCapacity=%d bytes, dataSize=%d bytes", maxCapacity, len(data))
	}

	// Compress data if required
	var indata []byte = data
	if defaultCompression {
		compressedData, err := c.CompressZSTD(data)
		if err != nil {
			return fmt.Errorf("failed to compress data: %w", err)
		}
		indata = compressedData
	}

	// Embed data
	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return fmt.Errorf("failed to embed data into RGB channels: %w", err)
	}

	// Generate image from embedded RGB channels
	imgdata, err := u.SaveImage(embeddedRGBChannels, height, width)
	if err != nil {
		return fmt.Errorf("failed to create embedded image: %w", err)
	}

	// Use default filename if none provided
	if outputFilename == "" {
		outputFilename = DefaultOutputFile
	}

	return SaveImage(outputFilename, imgdata)
}

// Decode extracts data embedded in an image using the specified bit depth.
// If the embedded data was compressed, it will be decompressed when `isDefaultCompressed` is true.
// Returns the extracted data or an error if the extraction or decompression fails.
//
// Parameters:
// - coverImage: The image containing embedded data to be extracted.
// - bitDepth: The bit depth used during the embedding process.
// - isDefaultCompressed: A flag indicating whether the embedded data was compressed.
func (m *ExtractHandler) Decode(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	// Validate coverImage dimensions
	if coverImage == nil {
		return nil, ErrInvalidCoverImage
	}
	if bitDepth < 0 || bitDepth > 7 {
		return nil, ErrDepthOutOfRange
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}
	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return nil, ErrFailedToExtractRGB
	}

	// Extract data
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, ErrFailedToExtractData
	}

	// Validate extracted data length
	lenData, err := u.GetlenOfData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get length of extracted data: %w", err)
	}

	if lenData == 0 {
		return nil, ErrInvalidDataLength
	}

	var moddedData = make([]byte, 0, lenData)
	defer func() {
		if r := recover(); r != nil {
			moddedData = nil
			err = fmt.Errorf("fatal error: %v", r)
		}
	}()

	for i := 4; i < lenData+4; i++ {
		if i >= len(data) {
			return nil, fmt.Errorf("index out of range while accessing data: %d", i)
		}
		moddedData = append(moddedData, data[i])
	}

	// Decompress data if required
	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress extracted data: %w", err)
		}
		return outdata, nil
	}

	return moddedData, nil
}

// Encode embeds data into a cover image using a specified bit depth, encrypts and compresses the data, and saves the resulting image to the specified output file.
// Parameters:
// - coverImage: The image to embed data into.
// - data: The data to embed in the image.
// - bitDepth: The bit depth used for embedding (valid range: 0-7).
// - outputFilename: The file name to save the resulting image. Defaults to a pre-defined name if empty.
// - password: The password used to encrypt the data.
// Returns:
// - error: An error if any part of the embedding process fails.
func (m *SecureEmbedHandler) Encode(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, password string) error {
	// Validate coverImage dimensions
	if coverImage == nil {
		return ErrInvalidCoverImage

	}

	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()
	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	// Validate bit depth
	if bitDepth < 0 || bitDepth > 7 {
		return ErrDepthOutOfRange
	}

	// Validate data
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}
	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return ErrFailedToExtractRGB
	}

	maxCapacity := (len(RGBchannels) * 3 * (int(bitDepth) + 1)) / 8
	if (len(data)*8)+32 > maxCapacity {
		return fmt.Errorf("data is too large to embed into the image: maxCapacity=%d bytes, dataSize=%d bytes", maxCapacity, len(data))
	}

	cipher, err := EncryptData(data, password)
	if err != nil {
		return err
	}

	compressedData, err := c.CompressZSTD(cipher)
	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	// Embed data
	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, compressedData, bitDepth)
	if err != nil {
		return fmt.Errorf("failed to embed data into RGB channels: %w", err)
	}

	// Generate image from embedded RGB channels
	imgdata, err := u.SaveImage(embeddedRGBChannels, height, width)
	if err != nil {
		return fmt.Errorf("failed to create embedded image: %w", err)
	}

	// Use default filename if none provided
	if outputFilename == "" {
		outputFilename = DefaultOutputFile
	}

	return SaveImage(outputFilename, imgdata)
}

// Decode extracts embedded data from a cover image using a specified bit depth, decrypts and decompresses it, and returns the original data.
// Parameters:
// - coverImage: The image containing the embedded data.
// - bitDepth: The bit depth used for extracting data (valid range: 0-7).
// - password: The password used to decrypt the embedded data.
// Returns:
// - []byte: The extracted original data.
// - error: An error if the extraction process fails.
func (m *SecureExtractHandler) Decode(coverImage image.Image, bitDepth uint8, password string) ([]byte, error) {
	// Validate coverImage dimensions
	if coverImage == nil {
		return nil, ErrInvalidCoverImage
	}

	if bitDepth < 0 || bitDepth > 7 {
		return nil, ErrDepthOutOfRange
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}
	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return nil, ErrFailedToExtractRGB
	}

	// Extract data
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, ErrFailedToExtractData
	}

	// Validate extracted data length
	lenData, err := u.GetlenOfData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get length of extracted data: %w", err)
	}
	if lenData == 0 {
		return nil, errors.New("extracted data length is zero")
	}

	var moddedData = make([]byte, 0, lenData)
	defer func() {
		if r := recover(); r != nil {
			moddedData = nil
			err = fmt.Errorf("fatal error: %v", r)
		}
	}()

	for i := 4; i < lenData+4; i++ {
		if i >= len(data) {
			return nil, fmt.Errorf("index out of range while accessing data: %d", i)
		}
		moddedData = append(moddedData, data[i])
	}

	outdata, err := c.DecompressZSTD(moddedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress extracted data: %w", err)
	}

	return DecryptData(outdata, password)
}

func openFiles(coverImagePath, dataFilePath string) (coverImage image.Image, dataFile []byte, err error) {
	cimg, err := Decodeimage(coverImagePath)
	if err != nil {
		return nil, nil, err
	}

	df, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, nil, err
	}

	return cimg, df, nil
}


// EmbedFile embeds data from a file into an image using a default bit depth of 1 (the last two bits in a byte).
// The data is first compressed and encrypted with the provided password before embedding into the image.
// Returns an error if the process fails at any stage.
//
// Parameters:
// - coverImagePath: The file path of the image to embed data into.
// - dataFilePath: The file path of the data to embed.
// - outputFilePath: The file path to save the resulting image with embedded data.
// - password: A password used to encrypt the data before embedding.

// ExtractFile extracts embedded data from an image using a default bit depth of 1 (the last two bits in a byte).
// The embedded data is decrypted and decompressed using the provided password.
// The extracted file is saved using its original name (stored within the embedded data).
// Returns an error if the process fails at any stage.
//
// Parameters:
// - coverImagePath: The file path of the image containing embedded data.
// - password: A password used to decrypt the embedded data after extraction.
func EmbedFile(coverImagePath, dataFilePath, outputFilePath, password string, bitDepth uint8) error {
	if coverImagePath == "" {
		return errors.New("invalid coverImagePath")
	}

	if dataFilePath == "" {
		return errors.New("invalid dataFilePath")
	}

	if outputFilePath == "" {
		return errors.New("invalid outputFilePath")
	}

	if password == "" {
		return errors.New("invalid password")
	}

	if bitDepth > 7 {
		return errors.New("invalid bit Depth")
	} 

	fp := filepath.Base(dataFilePath)
	ext := fmt.Sprintf("/-%s-/\n", fp)

	cf, df, err := openFiles(coverImagePath, dataFilePath)
	if err != nil {
		return err
	}

	df = append([]byte(ext), df...)

	var (
		wg       sync.WaitGroup
		erchan   = make(chan error)
		channels []u.RgbChannel
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		channels = u.ExtractRGBChannelsFromImageWithConCurrency(cf, runtime.NumCPU())
		if (len(df)*8)+32 > len(channels)*3*(int(bitDepth)+1) {
			erchan <- fmt.Errorf("error: Data too large to embed into the image")
			return
		}
	}()

	go func() {
		defer wg.Done()
		df, err = c.CompressZSTD(df)
		if err != nil {
			erchan <- err
			return
		}

		df, err = u.Encrypt(password, df)
		if err != nil {
			erchan <- err
			return
		}
	}()

	select {
	case <-erchan:
		return err
	default:
	}

	wg.Wait()

	channels, err = u.EmbedIntoRGBchannelsWithDepth(channels, df, bitDepth)
	if err != nil {
		return err
	}

	newImage, err := u.SaveImage(channels, cf.Bounds().Max.Y, cf.Bounds().Max.X)
	if err != nil {
		return nil
	}

	return SaveImage(outputFilePath, newImage)
}

// ExtractFile extracts embedded data from an image using a default bit depth of 1 (the last two bits in a byte).
// The embedded data is decrypted and decompressed using the provided password.
// The extracted file is saved using its original name (stored within the embedded data).
// Returns an error if the process fails at any stage.
//
// Parameters:
// - coverImagePath: The file path of the image containing embedded data.
// - password: A password used to decrypt the embedded data after extraction.
func ExtractFile(coverImagePath, password string, bitDepth uint8) error {
	if coverImagePath == "" {
		return errors.New("invalid coverImagePath")
	}

	if password == "" {
		return errors.New("invalid password")
	}

	if bitDepth > 7 {
		return errors.New("invalid bit Depth")
	} 

	cf, err := Decodeimage(coverImagePath)
	if err != nil {
		return err
	}

	channels := u.ExtractRGBChannelsFromImageWithConCurrency(cf, runtime.NumCPU())
	embeddedData, err := u.ExtractDataFromRGBchannelsWithDepth(channels, bitDepth)
	if err != nil {
		return err
	}

	lenData, err := u.GetlenOfData(embeddedData)
	if err != nil {
		return err
	}

	var cipherText = make([]byte, 0, lenData)
	defer func() {
		if r := recover(); r != nil {
			cipherText = nil
			err = fmt.Errorf("fatal error: %v", r)
		}
	}()

	for i := 4; i < lenData+4; i++ {
		if i >= len(embeddedData) {
			return fmt.Errorf("index out of range while accessing data: %d", i)
		}
		cipherText = append(cipherText, embeddedData[i])
	}

	plaintext, err := u.Decrypt(password, cipherText)
	if err != nil {
		return err
	}

	plaintext, err = c.DecompressZSTD(plaintext)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(plaintext))

	var filename string
	var size int
	if scanner.Scan() {
		filename = scanner.Text()
		size = len(scanner.Bytes())
	}

	ff, err := os.Create(filename[2 : len(filename)-2])
	if err != nil {
		return err
	}
	defer ff.Close()

	_, err = ff.Write(plaintext[size+1:])
	if err != nil {
		return err
	}

	return nil
}
