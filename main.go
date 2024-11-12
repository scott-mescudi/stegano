package main

import (
	"fmt"
	"image/png"
	"log"
	s "lsb/stegano"
	"os"
	"time"
)

func main() {
	inputfile := "testimages/in/input.png"
	outputfile := "testimages/out/out.png"

	file, err := os.Open(inputfile)
	if err!= nil {
        fmt.Println("Error opening file:", err)
        return
    }
	defer file.Close()

	data, err := os.ReadFile("main.txt")
	if err != nil {
		log.Fatal(err)
	}

	coverimage , err  := png.Decode(file)
	embedder := s.NewPngEncoder()

	fmt.Printf("Image can hold %d bytes\n", embedder.GetImageCapacity(coverimage))

	start := time.Now()
	err = embedder.EncodePngImage(coverimage, data, outputfile)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("total: ", time.Since(start))
		

	file2, err := os.Open(outputfile)
	if err!= nil {
        fmt.Println("Error opening file:", err)
        return
    }

	imagez , err  := png.Decode(file2)
	// Decode the embedded data from the image
	_, err = embedder.DecodePngImage(imagez)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// fmt.Println(string(embeddedData))
}

//TODO: implement huffman encoding for embedded data