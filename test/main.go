package main

import (
    "image/jpeg"
    "image/png"
    "os"
    "log"
)

func main() {
    // Open the JPEG file
    jpegFile, err := os.Open("../testimages/in/in.jpeg")
    if err != nil {
        log.Fatal(err)
    }
    defer jpegFile.Close()

    // Decode the JPEG image
    img, err := jpeg.Decode(jpegFile)
    if err != nil {
        log.Fatal(err)
    }

    // Create the PNG file to save the result
    pngFile, err := os.Create("output.png")
    if err != nil {
        log.Fatal(err)
    }
    defer pngFile.Close()

    // Encode the image to PNG format and save it to the new file
    err = png.Encode(pngFile, img)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("JPEG image has been successfully converted to PNG and saved as 'output.png'")
}
