# Stegano: A Steganography Library in Go
[![Tests](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml/badge.svg?event=push)](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml)
![GitHub License](https://img.shields.io/github/license/scott-mescudi/stegano)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/scott-mescudi/stegano)
[![Go Reference](https://pkg.go.dev/badge/github.com/scott-mescudi/stegano.svg)](https://pkg.go.dev/github.com/scott-mescudi/stegano)

Stegano is a Go library that provides tools for embedding and extracting data within images using steganographic techniques. The library currently supports PNG and JPEG image formats and includes ZSTD compression to optimize data storage within the images. Future improvements may include additional compression techniques like Huffman encoding.

## Features


- **Multi-Image Support:** Supports all images compatible with the `image.Image` type in Go.  
- **Data Compression:** Utilizes ZSTD compression for efficient embedding.  
- **Capacity Calculation:** Calculates the maximum data capacity of an image for embedding.  
- **Variable Depth Encoding:** Embeds bits up to and including the specified index.  
- **Concurrency:** Increases speed at the cost of memory usage.  
- **Embed at Certain Depth:** Embeds 1 bit per channel at the specified index.  
- **Save Image:** Efficient PNG encoding.  
- **Decode Image:** Helper function to decode images into the `image.Image` type.  

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
    "github.com/scott-mescudi/stegano"
)
```

### Working with Images

#### Encode Data into a Image

```go
package main

import (
	"log"

	"github.com/scott-mescudi/stegano"
)

func main() {
	coverImage, err := stegano.Decodeimage("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	encoder := stegano.NewEmbedHandler()
	
	err = encoder.EncodeAndSave(coverImage, []byte("Hello, World!"), stegano.MinBitDepth, stegano.DefaultpngOutputFile, true)     
	if err != nil {
		log.Fatalln(err)
	}
}
```

#### Decode Data from a Image

```go
package main

import (
	"fmt"
	"log"

	"github.com/scott-mescudi/stegano"
)

func main() {
	coverImage, err := stegano.Decodeimage("imageWithData.png")
	if err != nil {
		log.Fatalln(err)
	}

	encoder := stegano.NewExtractHandler()
	
	data, err := encoder.Decode(coverImage, stegano.MinBitDepth, true)     
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(data))
}
```


>> **⚠ Disclaimer:**  
>> - The image used for embedding data must always be a **.png** file or another format with no compression; otherwise, data will be lost.  
>> - The **bitDepth** must match the one originally used, or a slice-out-of-bounds error will occur.


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
- **Audio Support:** Add support for audio formats.
- **Multi carrier Support:** if muliple carrier files are provided split data into each carrier, also include what part of the data it is