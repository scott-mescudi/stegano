# Stegano: A Steganography Library in Go

[![Tests](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml/badge.svg?event=push)](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml)
![GitHub License](https://img.shields.io/github/license/scott-mescudi/stegano)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/scott-mescudi/stegano)
[![Go Reference](https://pkg.go.dev/badge/github.com/scott-mescudi/stegano.svg)](https://pkg.go.dev/github.com/scott-mescudi/stegano)

Stegano is a powerful, flexible Go library designed for securely embedding and extracting data within image files using advanced steganographic techniques. It allows developers to hide secret information within image files (such as PNG and JPEG formats) without visible changes to the image. This functionality is enabled through the manipulation of pixel data at varying depths, ensuring that data can be embedded or extracted effectively.

In addition to providing a seamless interface for image-based steganography, Stegano offers data compression support via ZSTD to optimize storage within images. This makes it a useful tool for embedding larger datasets while minimizing the file size. Future versions will expand support for additional compression methods and possibly other types of carrier files, such as audio.

---

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

To use this library, install it and its dependencies. Add the package to your Go project by running the following command:

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

---

## Working with Images

### Basic Usage

#### Encode Data into an Image

```go
package main

import (
    "log"
    "github.com/scott-mescudi/stegano"
)

func main() {
    coverImage, err := stegano.DecodeImage("image.png")
    if err != nil {
        log.Fatalln(err)
    }

    encoder := stegano.NewEmbedHandler()

    err = encoder.EncodeAndSave(coverImage, []byte("Hello, World!"), stegano.MinBitDepth, stegano.DefaultPngOutputFile, true)
    if err != nil {
        log.Fatalln(err)
    }
}
```

> **⚠ Disclaimer:**  
> - The image used for embedding data must always be a **.png** file or another format with no compression; otherwise, data will be lost.  

#### Decode Data from an Image

```go
package main

import (
    "fmt"
    "log"
    "github.com/scott-mescudi/stegano"
)

func main() {
    coverImage, err := stegano.DecodeImage("imageWithData.png")
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

> **⚠ Disclaimer:**  
> - The image used for extracting data must always be a **.png** file or another format with no compression; otherwise, data will be lost.  
> - The **bitDepth** must match the one originally used, or a slice-out-of-bounds error will occur.  

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
- **`(JpegEmbedder) EncodeJpegImage(coverImage image.Image, data []byte, outputFilename string) error`**: Embeds the provided data into the image and saves it to the specified output file.
- **`(JpegEmbedder) DecodeJpegImage(coverImage image.Image) ([]byte, error)`**: Extracts embedded data from a JPEG image.

---

## Benchmarks

#### Benchmark Results embedding (AMD Ryzen 5 7600X 6-Core Processor)

| **Library Name**              | **Test 1 (ns/op)** | **Test 2 (ns/op)** | **Test 3 (ns/op)** | **Avg Time (ns/op)** | **Avg Time (ms/op)** |
|-------------------------------|--------------------|--------------------|--------------------|----------------------|----------------------|
| **Stegano**                    | `352,531,567`      | `349,444,200`      | `348,196,967`      | **`350,390,578`**    | **`350.39 ms`**      |
| **Stegano with concurrency**   | `286,168,125`      | `293,260,925`      | `284,079,175`      | **`287,169,408`**    | **`287.17 ms`**      |
| [**auyer/steganography**](https://github.com/auyer/steganography)       | `1,405,256,700`    | `1,424,957,200`    | `1,401,682,600`    | **`1,410,965,500`**  | **`1,410.97 ms`**    |

> **Image size:** 10,473,459 bytes  
> **Text size:** 641,788 bytes  
> Benchmark code can be found [here](./examples/steganobench)

---

## Future Improvements

- **Huffman Encoding:** Add support for Huffman-based compression.
- **Audio Support:** Add support for audio formats.
- **Multi-Carrier Support:** If multiple carrier files are provided, split data into each carrier and include what part of the data it contains.
