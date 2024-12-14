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
func (m *embedHandler) Encode(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	// Validate coverImage dimensions
	if coverImage == nil {
		return errors.New("coverImage is nil")
	}
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()
	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	// Validate bit depth
	if bitDepth < 0 || bitDepth > 7 {
		return fmt.Errorf("bitDepth is out of range (0-7): %d", bitDepth)
	}

	// Validate data
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return errors.New("failed to extract RGB channels from the image")
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
		outputFilename = DefaultpngOutputFile
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
func (m *extractHandler) Decode(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	// Validate coverImage dimensions
	if coverImage == nil {
		return nil, errors.New("coverImage is nil")
	}
	if bitDepth < 0 || bitDepth > 7 {
		return nil, fmt.Errorf("bitDepth is out of range (1-7): %d", bitDepth)
	}

	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return nil, errors.New("failed to extract RGB channels from the image")
	}

	// Extract data
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from RGB channels: %w", err)
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



func (m *secureEmbedHandler) Encode(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, password string) error {
	// Validate coverImage dimensions
	if coverImage == nil {
		return errors.New("coverImage is nil")

	}

	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()
	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	// Validate bit depth
	if bitDepth < 0 || bitDepth > 7 {
		return fmt.Errorf("bitDepth is out of range (0-7): %d", bitDepth)
	}

	// Validate data
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return errors.New("failed to extract RGB channels from the image")
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
		outputFilename = DefaultpngOutputFile
	}

	
	return SaveImage(outputFilename, imgdata)
}

func (m *secureExtractHandler) Decode(coverImage image.Image, bitDepth uint8, password string) ([]byte, error) {
	// Validate coverImage dimensions
	if coverImage == nil {
		return nil, errors.New("coverImage is nil")
	}

	if bitDepth < 0 || bitDepth > 7 {
		return nil, fmt.Errorf("bitDepth is out of range (1-7): %d", bitDepth)
	}

	// Extract RGB channels
	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if RGBchannels == nil {
		return nil, errors.New("failed to extract RGB channels from the image")
	}

	// Extract data
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from RGB channels: %w", err)
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

const bd = 1

func EmbedFile(coverImagePath, dataFilePath, outputFilePath, password string) (error) {
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

	fp := filepath.Base(dataFilePath)
	ext := fmt.Sprintf("/-%s-/\n", fp)

	cf, df, err := openFiles(coverImagePath, dataFilePath)
	if err != nil {
		return err
	}

	df = append([]byte(ext), df...)

	// conc
	df, err = c.CompressZSTD(df)
	if err != nil {
		return err
	}

	// conc
	df, err = u.Encrypt(password, df)
	if err != nil {
		return err
	}

	// conc
	channels := u.ExtractRGBChannelsFromImageWithConCurrency(cf, runtime.NumCPU())
	if (len(df)*8)+32 > len(channels)*3*(int(bd)+1) {
		return fmt.Errorf("error: Data too large to embed into the image")
	}

	channels, err = u.EmbedIntoRGBchannelsWithDepth(channels, df, bd)
	if err != nil {
		return err
	}


	newImage, err := u.SaveImage(channels, cf.Bounds().Max.Y, cf.Bounds().Max.X)
	if err != nil {
		return nil
	}

	return SaveImage(outputFilePath, newImage)
}

func ExtractFile(coverImagePath, password string) (error) {
	if coverImagePath == "" {
		return errors.New("invalid coverImagePath")
	}

	if password == "" {
		return errors.New("invalid password")
	}

	cf, err := Decodeimage(coverImagePath)
	if err != nil {
		return err
	}


	channels := u.ExtractRGBChannelsFromImageWithConCurrency(cf, runtime.NumCPU())
	embeddedData, err := u.ExtractDataFromRGBchannelsWithDepth(channels, bd)
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
            return  fmt.Errorf("index out of range while accessing data: %d", i)
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


	ff, err := os.Create(filename[2:len(filename)-2])
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