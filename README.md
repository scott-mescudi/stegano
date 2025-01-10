# Stegano: The fastest Steganography Library for Go

[![Tests](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml/badge.svg?event=push)](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml) ![GitHub License](https://img.shields.io/github/license/scott-mescudi/stegano) [![Go Reference](https://pkg.go.dev/badge/github.com/scott-mescudi/stegano.svg)](https://pkg.go.dev/github.com/scott-mescudi/stegano)  [![Go Report Card](https://goreportcard.com/badge/github.com/scott-mescudi/stegano)](https://goreportcard.com/badge/github.com/scott-mescudi/stegano)

## Table of Contents

1. [Features](#features)
2. [What is Steganography?](#what-is-steganography)
3. [Use Cases](./docs/usecases.md)
4. [Installation](#installation)
5. [Usage](#usage)
    - [Import the Library](#import-the-library)
6. [Working with Images](#working-with-images)
	- [Quickstart](#Quickstart)
    - [Embed a Message into an Image](#1-embed-a-message-into-an-image)
    - [Extract a Message from an Embedded Image](#2-extract-a-message-from-an-embedded-image)
    - [Embed Data Without Compression](#3-embed-data-without-compression)
    - [Extract Data Without Compression](#4-extract-data-without-compression)
    - [Embed at a Specific Bit Depth](#5-embed-at-a-specific-bit-depth)
    - [Extract Data from a Specific Bit Depth](#6-extract-data-from-a-specific-bit-depth)
    - [Check Image Capacity](#7-check-image-capacity)
    - [Embed Encrypted Data](#8-embed-encrypted-data)
    - [Extract and Decrypt Data](#9-extract-and-decrypt-data)
7. [Working with Audio](#Working-with-Audio-(Experimental))
    - [Embed Data into WAV Files](#1-embed-data-into-wav-files)
    - [Extract Data from WAV Files](#2-extract-data-from-wav-files)
    - [Embed at Specific Bit Depth](#3-embed-at-specific-bit-depth)
    - [Extract from Specific Bit Depth](#4-extract-from-specific-bit-depth)
8. [Advanced Options](#advanced-options)
9. [Notes](#notes)
10. [Benchmarks](#benchmarks)
11. [Future Improvements](#future-improvements)

---

## Features

- **Multi-Image Support**: Works with any image type compatible with Go's `image.Image`.
- **Data Compression**: Supports ZSTD compression to minimize the size of embedded data.
- **Reed-Solomon Codes**: Implements Reed-Solomon error correction.
- **Capacity Calculation**: Automatically calculates the maximum capacity of an image for data embedding.
- **Variable Depth Encoding**: Allows you to embed data up to a specified bit depth.
- **Concurrency**: Supports concurrent processing for improved speed (higher memory usage).
- **Custom Bit Depth Embedding**: Lets you specify the bit depth used for data embedding (e.g., LSB, MSB).
- **Encryption**: Enables secure encryption of data before embedding into the image.
- **Efficient PNG Encoding**: Saves the image in PNG format.

---

## What is Steganography?

Steganography is the practice of concealing data within other, seemingly innocent data in such a way that it remains undetectable to the casual observer. Unlike encryption, which focuses on making data unreadable, steganography hides it entirely, making the data appear normal and undisturbed.

The most common use case for steganography is hiding text or binary data within an image. This library enables you to easily embed and extract such hidden information from images.

---

## Installation

To integrate Stegano into your Go project, simply run:

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

### Quickstart

The following are high-level functions for embedding and extracting data easily:

```go
func main() {
	err := stegano.EmbedFile("cover.png", "data.txt", stegano.DefaultOutputFile, "password123", stegano.LSB)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
```

```go
func main() {
	err := stegano.ExtractFile(stegano.DefaultOutputFile, "password123", stegano.LSB)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
```

For more control over the process, refer to the examples below:

### 1. Embed a Message into an Image

You can embed a message into an image using the `EmbedHandler` class.

```go
func main() {
	// wrapper function around different image decoders.
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	embedder := stegano.NewEmbedHandler()

	// Encode and save the message in the cover image.
	err = embedder.Encode(coverFile, []byte("Hello, World!"), stegano.MaxBitDepth, stegano.DefaultOutputFile, true)
	if err != nil {
		log.Fatalln(err)
	}
}
```

### 2. Extract a Message from an Embedded Image

Extract a hidden message from an image using the `ExtractHandler` class.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	extractor := stegano.NewExtractHandler()

	// Decode the message from the image.
	data, err := extractor.Decode(coverFile, stegano.MaxBitDepth, true)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

### 3. Embed Data Without Compression

Embed data into an image without using any compression.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	embedder := stegano.NewEmbedHandler()

	// Embed the message into the image without compression.
	embeddedImage, err := embedder.EmbedDataIntoImage(coverFile, []byte("Hello, World!"), stegano.LSB)
	if err != nil {
		log.Fatalln(err)
	}

	err = stegano.SaveImage(stegano.DefaultOutputFile, embeddedImage)
	if err != nil {
		log.Fatalln(err)
	}
}
```

### 4. Extract Data Without Compression

Extract data from an image where no compression was used.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	extractor := stegano.NewExtractHandler()

	// Extract uncompressed data from the image.
	data, err := extractor.ExtractDataFromImage(coverFile, stegano.LSB)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

### 5. Embed at a Specific Bit Depth

Embed data at a specific bit depth for better control over the hiding technique.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	embedder := stegano.NewEmbedHandler()

	// Embed the message at a specific bit depth (e.g., 3).
	embeddedImage, err := embedder.EmbedAtDepth(coverFile, []byte("Hello, World!"), 3)
	if err != nil {
		log.Fatalln(err)
	}

	err = stegano.SaveImage(stegano.DefaultOutputFile, embeddedImage)
	if err != nil {
		log.Fatalln(err)
	}
}
```

### 6. Extract Data from a Specific Bit Depth

Extract data from an image using the same bit depth used for embedding.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	extractor := stegano.NewExtractHandler()

	// Extract data from the image at a specific bit depth (e.g., 3).
	data, err := extractor.ExtractAtDepth(coverFile, 3)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

### 7. Check Image Capacity

Check how much data an image can hold based on the bit depth.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Calculate and print the data capacity of the image.
	capacity := stegano.GetImageCapacity(coverFile, stegano.MaxBitDepth)
	fmt.Printf("Image capacity at bit depth %d: %d bytes\n", stegano.MaxBitDepth, capacity)
}
```

### 8. Embed Encrypted Data

Encrypt the data before embedding it into the image for added security.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the data before embedding.
	encryptedData, err := stegano.EncryptData([]byte("Hello, World!"), "your-encryption-key")
	if err != nil {
		log.Fatalln(err)
	}

	embedder := stegano.NewEmbedHandler()

	// Embed the encrypted data into the image.
	err = embedder.Encode(coverFile, encryptedData, stegano.LSB, stegano.DefaultOutputFile, true)
	if err != nil {
		log.Fatalln(err)
	}
}
```

### 9. Extract and Decrypt Data

Extract encrypted data from the image and decrypt it.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	extractor := stegano.NewExtractHandler()

	// Extract the encrypted data.
	encryptedData, err := extractor.Decode(coverFile, stegano.LSB, true)
	if err != nil {
		log.Fatalln(err)
	}

	// Decrypt the data.
	decryptedData, err := stegano.DecryptData(encryptedData, "your-encryption-key")
	if err != nil {
		log.Fatalln(err)
	}

	// Print the decrypted message.
	fmt.Println(string(decryptedData))
}
```

---

## Working with Audio (Experimental)

### 1. Embed Data into WAV Files

```go
func main() {
    embedder := stegano.NewAudioEmbedHandler()
    
    err := embedder.EmbedIntoWAVWithDepth("input.wav", "output.wav", []byte("Hello World"), stegano.LSB)
    if err != nil {
        log.Fatalln(err)
    }
}
```

### 2. Extract Data from WAV Files

```go
func main() {
    extractor := stegano.NewAudioExtractHandler()
    
    data, err := extractor.ExtractFromWAVWithDepth("embedded.wav", stegano.LSB)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println(string(data))
}
```

### 3. Embed at Specific Bit Depth

```go
func main() {
    embedder := stegano.NewAudioEmbedHandler()
    
    err := embedder.EmbedIntoWAVAtDepth("input.wav", "output.wav", []byte("Hello World"), 3)
    if err != nil {
        log.Fatalln(err)
    }
}
```

### 4. Extract from Specific Bit Depth

```go
func main() {
    extractor := stegano.NewAudioExtractHandler()
    
    data, err := extractor.ExtractFromWAVAtDepth("embedded.wav", 3)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println(string(data))
}
```

---

## Advanced Options

> You can use concurrency to speed up the embedding and extraction process. Use the `NewEmbedHandlerWithConcurrency` or `NewExtractHandlerWithConcurrency` functions to specify the number of goroutines to be used.

```go
embedder, err := stegano.NewEmbedHandlerWithConcurrency(12)
if err != nil {
	log.Fatalln(err)
}

extractor, err := stegano.NewExtractHandlerWithConcurrency(12)
if err != nil {
	log.Fatalln(err)
}
```

---

## Notes

> - This library can be used with any image type but works best with **PNG** images.
> - **Bit Depth and Compression**: Ensure that the same bit depth and compression settings are used during both embedding and extraction.
> - **Default Output Format**: By default, images NEED to be saved in PNG format to avoid any data loss.

---

## Benchmarks

| **Library Name**               | **Test 1 (ns/op)** | **Test 2 (ns/op)** | **Test 3 (ns/op)** | **Avg Time (ns/op)** | **Avg Time (ms/op)** |
|---------------------------------|--------------------|--------------------|--------------------|----------------------|----------------------|
| **Stegano**                     | 352,531,567        | 349,444,200        | 348,196,967        | **350,390,578**      | **350.39 ms**        |
| **Stegano with Concurrency**    | 286,168,125        | 293,260,925        | 284,079,175        | **287,169,408**      | **287.17 ms**        |
| [**auyer/steganography**](https://github.com/auyer/steganography) | 1,405,256,700      | 1,424,957,200      | 1,401,682,600      | **1,410,965,500**    | **1,410.97 ms**      |

> **Image size:** 10,473,459 bytes  
> **Text size:** 641,788 bytes  
> Benchmark code can be found [here](./examples/steganobench)

---

## Future Improvements
> - **More tests**
> - **Huffman Encoding**: Add support for more efficient compression techniques like Huffman coding.
> - **Multi-Carrier Support**: Enable splitting data across multiple images or files for larger data embedding.