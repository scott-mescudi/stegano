package main

import (
	"fmt"
	"image/png"
	s "lsb/internal"
	"os"
)

func main() {
	inputfile := "D:\\lsb\\testimages\\in\\input.png"
	outputfile := "D:\\lsb\\testimages\\out\\output.png"

	file, err := os.Open(inputfile)
	if err!= nil {
        fmt.Println("Error opening file:", err)
        return
    }

	imagev , err  := png.Decode(file)

	embedder := s.NewPngEncoder()

	// Encode data into the image
	err = embedder.EncodePngImage(imagev, []byte("skibdidi"), outputfile)
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
}

//TODO: implement huffman encoding for embedded data