package main

import (
	"fmt"
	"image/png"
	"log"
	s "lsb/stegano"
	"os"
)

func main() {
	inputfile := "testimages/in/input.png"
	outputfile := "testimages/out/out.png"

	file, err := os.Open(inputfile)
	if err!= nil {
        fmt.Println("Error opening file:", err)
        return
    }

	data, err := os.ReadFile("main.txt")
	if err != nil {
		log.Fatal(err)
	}

	coverimage , err  := png.Decode(file)
	embedder := s.NewPngEncoder()

	fmt.Printf("Image can hold %d bytes", embedder.GetImageCapacity(coverimage))

	// Encode data into the image
	err = embedder.EncodePngImage(coverimage, data, outputfile)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
		

	file2, err := os.Open(outputfile)
	if err!= nil {
        fmt.Println("Error opening file:", err)
        return
    }

	imagez , err  := png.Decode(file2)

	// Decode the embedded data from the image
	embeddedData, err := embedder.DecodePngImage(imagez)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	fmt.Println(string(embeddedData))
}//

//TODO: implement huffman encoding for embedded data