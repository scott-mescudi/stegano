package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/auyer/steganography"
	"github.com/scott-mescudi/stegano"
)

func steganoEmbed(filename string, data []byte) error {
	steg := stegano.NewEmbedHandler()

	img, err := stegano.Decodeimage(filename)
	if err != nil {
		return err
	}


	modimg, err := steg.EmbedDataIntoImage(img, data, 0)
	if err != nil {
		return err
	}


	if err := stegano.SaveImage(stegano.DefaultpngOutputFile, modimg); err != nil {
		return err
	}

	return nil
}

func steganoExtract(filename string) ([]byte, error) {
	steg := stegano.NewExtractHandler()

	img, err := stegano.Decodeimage(filename)
	if err != nil {
		return nil, err
	}

	data, err := steg.ExtractDataFromImage(img, 0)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func auyerEmbed(inputFilename, outputFilename string, data []byte) error {
	inFile, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %v", inputFilename, err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, err := png.Decode(reader)
	if err != nil {
		return fmt.Errorf("failed to decode PNG from file %s: %v", inputFilename, err)
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, data)
	if err != nil {
		return fmt.Errorf("failed to encode message into image: %v", err)
	}

	outFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFilename, err)
	}
	defer outFile.Close()

	_, err = w.WriteTo(outFile)
	if err != nil {
		return fmt.Errorf("failed to write encoded image to file %s: %v", outputFilename, err)
	}

	return nil
}

func auyerExtract(inputFilename string) ([]byte, error) {
	inFile, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file %s: %v", inputFilename, err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, err := png.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG from file %s: %v", inputFilename, err)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)
	msg := steganography.Decode(sizeOfMessage, img)
	if msg == nil {
		return nil, fmt.Errorf("failed to decode message from file %s", inputFilename)
	}

	return msg, nil
}

func main() {
	// Read data for embedding
	data, err := os.ReadFile("data.txt")
	if err != nil {
		log.Println("Error:", err)
		return
	}

	// Testing scott-mescudi/stegano
	fmt.Println("Testing scott-mescudi/stegano")
	start := time.Now()
	err = steganoEmbed("./in/big.png", data)
	if err != nil {
		log.Println("Error during stegano embed:", err)
		return
	}
	fmt.Printf("Stegano Embed took: %v\n", time.Since(start))

	start = time.Now()
	od, err := steganoExtract(stegano.DefaultpngOutputFile)
	if err != nil {
		log.Println("Error during stegano extract:", err)
		return
	}
	fmt.Printf("Stegano Extract took: %v\n", time.Since(start))
	fmt.Printf("Stegano: %v -> %v\n", len(data), len(od))

	// Testing auyer/steganography
	fmt.Println("\nTesting auyer/steganography")
	start = time.Now()
	err = auyerEmbed("./in/big.png", "out.png", data)
	if err != nil {
		log.Println("Error during auyer embed:", err)
		return
	}
	fmt.Printf("Auyer Embed took: %v\n", time.Since(start))

	start = time.Now()
	od, err = auyerExtract("out.png")
	if err != nil {
		log.Println("Error during auyer extract:", err)
		return
	}
	fmt.Printf("Auyer Extract took: %v\n", time.Since(start))
	fmt.Printf("Auyer: %v -> %v\n", len(data), len(od))
}