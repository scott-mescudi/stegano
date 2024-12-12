# Stegano: A Steganography Library in Go

[![Tests](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml/badge.svg?event=push)](https://github.com/scott-mescudi/stegano/actions/workflows/go.yml)
![GitHub License](https://img.shields.io/github/license/scott-mescudi/stegano) [![Go Reference](https://pkg.go.dev/badge/github.com/scott-mescudi/stegano.svg)](https://pkg.go.dev/github.com/scott-mescudi/stegano)


Stegano is a Go library that provides tools for embedding and extracting data within images using steganographic techniques. The library currently supports any image that support the image.Iamge type in go and includes ZSTD compression to optimize data storage within images.

---

## Table of Contents

1. [Features](#features)
2. [What is Steganography?](#what-is-steganography)
3. [Installation](#installation)
4. [Usage](#usage)
   - [Import the Library](#import-the-library)
5. [Working with Images](#working-with-images)
    - [Embed a Message into an Image](#1-embed-a-message-into-an-image)
    - [Extract a Message from an Embedded Image](#2-extract-a-message-from-an-embedded-image)
    - [Embed Data Without Compression](#3-embed-data-without-compression)
    - [Extract Data Without Compression](#4-extract-data-without-compression)
    - [Embed at a Specific Bit Depth](#5-embed-at-a-specific-bit-depth)
    - [Extract Data from a Specific Bit Depth](#6-extract-data-from-a-specific-bit-depth)
    - [Check Image Capacity](#7-check-image-capacity)
6. [Advanced Options](#advanced-options)
7. [Notes](#notes)
8. [Benchmarks](#benchmarks)
9. [Future Improvements](#future-improvements)

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

## What is Steganography?

Steganography is the practice of hiding data inside other, non-suspicious data in such a way that it is imperceptible to an observer. This is different from encryption, where data is scrambled to make it unreadable, but still detectable. In steganography, the goal is to hide the data so that it goes unnoticed.

### Example Process

Below is an example of how steganography works with an image. The original image, embedded data, and resulting image appear unchanged to the human eye.

| **Original Image** | **Embedded Data** | **Resulting Image** |
|--------------------|-------------------|---------------------|
| ![Original Image](./examples/assets/in.png) | `Hello, World!` (Encoded in LSBs) | ![Resulting Image](./examples/assets/out.png) |

In this example, the message "Hello, World!" is hidden within the image, but the image looks the same as the original one to the naked eye (try to extract it using the tool with bitDepth 0 and compression).

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

# Working with Images


### **1. Embed a Message into an Image**
This function demonstrates how to embed a message into a cover image.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an EmbedHandler instance for managing the embedding process.
	embedder := stegano.NewEmbedHandler()

        // Encode and save the message in the cover image.
        // The settings used are:
        // - Minimum bit depth for embedding.
        // - PNG output format by default.
        // - Optional compression enabled.

        // This library by default embeds data up to and including the index.
        // For example, if embedding at depth 3:
        // - A binary value of 0x11111111 will become 0x11110000.

        // to embed into the LSB use a bitdepth of 0 or stegano.MinBitDepth.

	err = embedder.EncodeAndSave(coverFile, []byte("Hello, World!"), stegano.MinBitDepth, stegano.DefaultpngOutputFile, true)
	if err != nil {
		log.Fatalln(err)
	}
}
```

> **Note:** You can also enable concurrency for embedding using `NewEmbedHandlerWithConcurrency(numGoroutines)` for better performance with large files.

---

### **2. Extract a Message from an Embedded Image**
This function demonstrates how to retrieve a hidden message from an image.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an ExtractHandler instance for managing the extraction process.
	extractor := stegano.NewExtractHandler()

	// Decode the message from the image.
	// Ensure the following settings match those used during embedding:
	// - Minimum bit depth.
	// - Optional compression enabled/disabled.
	data, err := extractor.Decode(coverFile, stegano.MinBitDepth, true)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

> **Disclaimer:** Ensure that the bit depth, compression settings, and other parameters used during embedding are known, as incorrect settings may result in corrupted or irretrievable data.

---

### **3. Embed Data Without Compression**
This method embeds data without the use of automatic compression.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an EmbedHandler instance for managing the embedding process.
	embedder := stegano.NewEmbedHandler()

	// Embed the message into the image without saving to a file directly.
	embeddedImage, err := embedder.EmbedDataIntoImage(coverFile, []byte("Hello, World!"), stegano.MinBitDepth)
	if err != nil {
		log.Fatalln(err)
	}

	// Save the modified image to a file.
	err = stegano.SaveImage(stegano.DefaultpngOutputFile, embeddedImage)
	if err != nil {
		log.Fatalln(err)
	}
}
```

---

### **4. Extract Data Without Compression**
This function demonstrates how to extract uncompressed data.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an ExtractHandler instance for managing the extraction process.
	extractor := stegano.NewExtractHandler()

	// Extract data from the image.
	data, err := extractor.ExtractDataFromImage(coverFile, stegano.MinBitDepth)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

---

### **5. Embed at a Specific Bit Depth**
For fine-grained control, embed data at a specific bit depth.

```go
func main() {
	coverFile, err := stegano.Decodeimage("coverimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an EmbedHandler instance for managing the embedding process.
	embedder := stegano.NewEmbedHandler()

	// Embed the message at a specific bit depth.
	// For example, if embedding at depth 3:
	// - A binary value of 0x11111111 will become 0x11110111.
	embeddedImage, err := embedder.EmbedAtDepth(coverFile, []byte("Hello, World!"), 3)
	if err != nil {
		log.Fatalln(err)
	}

	// Save the modified image to a file.
	err = stegano.SaveImage(stegano.DefaultpngOutputFile, embeddedImage)
	if err != nil {
		log.Fatalln(err)
	}
}
```

> **Disclaimer:** Altering specific bit depths can affect image quality. Use bit depth settings carefully to maintain a balance between data embedding and visual fidelity.

---

### **6. Extract Data from a Specific Bit Depth**
Retrieve data from an image, ensuring the correct bit depth is used.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Create an ExtractHandler instance for managing the extraction process.
	extractor := stegano.NewExtractHandler()

	// Extract data from the image at the specified bit depth.
	data, err := extractor.ExtractAtDepth(coverFile, 3)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the extracted message.
	fmt.Println(string(data))
}
```

> **Disclaimer:** Extraction must use the exact bit depth specified during embedding. Mismatched settings will likely result in errors or incorrect data.

---

### **7. Check Image Capacity**
Determine the maximum data capacity of an image for a given bit depth.

```go
func main() {
	coverFile, err := stegano.Decodeimage("embeddedimage.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Calculate and print the data capacity for the given image.
	capacity := stegano.GetImageCapacity(coverFile, stegano.MinBitDepth)
	fmt.Printf("Image capacity at bit depth %d: %d bytes\n", stegano.MinBitDepth, capacity)
}
```

> **Disclaimer:** The maximum data capacity depends on the image size and bit depth. Be aware that embedding too much data can visibly degrade the image.

---

## Advanced Options
The library provides concurrency options for embedding and extraction. Use `NewEmbedHandlerWithConcurrency` or `NewExtractHandlerWithConcurrency` to control the number of goroutines for faster processing.

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
- Always use the same settings (bit depth, compression) for embedding and extraction.
- Default output format MUST be PNG or any lossless image formats.

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
