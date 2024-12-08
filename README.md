# Stegano: A Steganography Library in Go
[![Tests](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml/badge.svg?event=push)](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml)
![GitHub License](https://img.shields.io/github/license/scott-mescudi/stegano)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/scott-mescudi/stegano)
[![Go Reference](https://pkg.go.dev/badge/github.com/scott-mescudi/stegano.svg)](https://pkg.go.dev/github.com/scott-mescudi/stegano)


Stegano is a Go library that provides tools for embedding and extracting data within images using steganographic techniques. The library currently supports PNG and JPEG image formats and includes ZSTD compression to optimize data storage within the images. Future improvements may include additional compression techniques like Huffman encoding.

## Features

- **PNG Image Support:** Embed and extract data from PNG images.
- **JPEG Image Support:** Embed and extract data from JPEG images.
- **Data Compression:** Utilizes ZSTD compression for efficient embedding.
- **Capacity Calculation:** Calculate the maximum data capacity of an image for embedding.
- **Variable Depth Encoding:** Can adjust how many bits in a byte are used to store data.


---

## Installation

To use this library, you'll need to install it and its dependencies. Add the package to your Go project by importing it and its required dependencies:

```bash
go get github.com/scott-mescudi/stegano@latest
```
---

## Usage

### Import the Library

```go
import (
    "image"
    "github.com/scott-mescudi/stegano"
)
```

### Working with PNG Images

#### Encode Data into a PNG Image

```go
package main

import (
    "image"
    "os"
    "github.com/scott-mescudi/stegano"
    _ "image/png" // Required for PNG decoding
)

func main() {
    file, _ := os.Open("cover.png")
    defer file.Close()
    
    img, _ := png.Decode(file)
    data := []byte("Secret message to hide")
    outputFile := "output.png"
    
    pngEmbedder := stegano.NewPngEncoder()
    err := pngEmbedder.EncodePngImage(img, data, outputFile)
    if err != nil {
        panic(err)
    }
}
```

#### Decode Data from a PNG Image

```go
package main

import (
    "image"
    "os"
    "github.com/scott-mescudi/stegano"
    _ "image/png" // Required for PNG decoding
)

func main() {
    file, _ := os.Open("output.png")
    defer file.Close()
    
    img, _ := png.Decode(file)
    
    pngEmbedder := stegano.NewPngEncoder()
    data, err := pngEmbedder.DecodePngImage(img)
    if err != nil {
        panic(err)
    }
    
    println(string(data)) // Output the hidden message
}
```

---

### Working with JPEG Images

#### Encode Data into a JPEG Image

```go
package main

import (
    "image"
    "os"
    "github.com/scott-mescudi/stegano"
    _ "image/jpeg" // Required for JPEG decoding
)

func main() {
    file, _ := os.Open("cover.jpg")
    defer file.Close()
    
    img, _ := jpeg.Decode(file)
    data := []byte("Hidden JPEG data")
    outputFile := "output.jpg"
    
    jpegEmbedder := stegano.NewJpegEncoder()
    err := jpegEmbedder.EncodeJPEGImage(img, data, outputFile)
    if err != nil {
        panic(err)
    }
}
```

#### Decode Data from a JPEG Image

```go
package main

import (
    "image"
    "os"
    "github.com/scott-mescudi/stegano"
    _ "image/jpeg" // Required for JPEG decoding
)

func main() {
    file, _ := os.Open("output.jpg")
    defer file.Close()
    
    img, _ := jpeg.Decode(file)
    
    jpegEmbedder := stegano.NewJpegEncoder()
    data, err := jpegEmbedder.DecodeJPEGImage(img)
    if err != nil {
        panic(err)
    }
    
    println(string(data)) // Output the hidden message
}
```

---

## API Documentation

### PNG Embedder

- **`NewPngEncoder() PngEmbedder`**: Returns a new instance of `PngEmbedder`.
- **`(PngEmbedder) GetImageCapacity(coverImage image.Image) int`**: Returns the maximum number of bytes that can be embedded in a PNG image.
- **`(PngEmbedder) EncodePngImage(coverImage image.Image, data []byte, outputFilename string) error`**: Embeds the provided data into the image and saves it to the specified output file.
- **`(PngEmbedder) DecodePngImage(coverImage image.Image) ([]byte, error)`**: Extracts embedded data from a PNG image.

### JPEG Embedder

- **`NewJpegEncoder() JpegEmbedder`**: Returns a new instance of `JpegEmbedder`.
- **`(JpegEmbedder) GetImageCapacity(coverImage image.Image) int`**: Returns the maximum number of bytes that can be embedded in a JPEG image.
- **`(JpegEmbedder) EncodeJPEGImage(coverImage image.Image, data []byte, outputFilename string) error`**: Embeds the provided data into the image and saves it to the specified output file.
- **`(JpegEmbedder) DecodeJPEGImage(coverImage image.Image) ([]byte, error)`**: Extracts embedded data from a JPEG image.

---

## Future Improvements

- **Huffman Encoding:** Add support for Huffman-based compression.
- **Encryption:** Integrate encryption options for better data security.
- **Audio Support:** Add support for audio formats.