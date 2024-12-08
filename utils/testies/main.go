package main

import (
	"fmt"
	s "github.com/scott-mescudi/stegano"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

func jpegtest(comp bool) {
	inputfile := "testimages/in/in.jpeg"
	outputfile := "testimages/out/outjpeg.png"

	file, err := os.Open(inputfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := os.ReadFile("main.txt")
	if err != nil {
		log.Fatal(err)
	}

	coverimage, err := jpeg.Decode(file)
	embedder := s.NewJpegEncoder()

	fmt.Printf("Image can hold %d bytes\n", embedder.GetImageCapacity(coverimage))

	start := time.Now()
	err = embedder.EncodeJPEGImage(coverimage, data, outputfile, comp)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("total: ", time.Since(start))

	file2, err := os.Open(outputfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	imagez, err := png.Decode(file2)
	_, err = embedder.DecodeJPEGImage(imagez, comp)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// fmt.Println(string(embeddedData))
}

func pngtest(comp bool) {
	embedder := s.NewPngEncoder()
	outputfile := "testimages/out/out.png"
	inputfile := "D:\\go\\lsb\\testies\\testimages\\in\\images.png"

	file, err := os.Open(inputfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := os.ReadFile("main.txt")
	if err != nil {
		log.Fatal(err)
	}

	coverimage, err := png.Decode(file)

	fmt.Printf("Image can hold %d bytes\n", embedder.GetImageCapacity(coverimage))

	start := time.Now()
	err = embedder.EncodePngImage(coverimage, data, outputfile, comp)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("total: ", time.Since(start))

	file2, err := os.Open(outputfile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	imagez, err := png.Decode(file2)
	// Decode the embedded data from the image
	embeddedData, err := embedder.DecodePngImage(imagez, comp)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	fmt.Println(string(embeddedData))
}

func hsdat() {
	pg := s.NewPngEncoder()

	of := "D:\\go\\lsb\\testies\\testimages\\in\\images.png"
	of2 := "testimages/out/out.png"

	file2, err := os.Open(of)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	file3, err := os.Open(of2)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	imagez, err := png.Decode(file2)
	imaget, err := png.Decode(file3)

	fmt.Println(pg.HasData(imagez))
	fmt.Println(pg.HasData(imaget))
}

func main() {
	pngtest(false)
	hsdat()
}
