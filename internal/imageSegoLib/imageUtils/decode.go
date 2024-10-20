package imageutils

import (
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "fmt"
)

// decodeImage decodes an image based on its file extension (supports .png, .jpg, .jpeg)
func DecodeImage(filename string) (image.Image, error) {
    ext := filepath.Ext(filename)
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("error opening image: %w", err)
    }
    defer file.Close()

    switch ext {
    case ".png":
        return png.Decode(file)
    case ".jpeg", ".jpg":
        return jpeg.Decode(file)
    default:
        return nil, fmt.Errorf("unsupported file format: %s", ext)
    }
}
