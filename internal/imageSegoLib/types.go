package imagesegolib

// Struct for holding RGB channel data as 32-bit
type RgbChannel struct {
    r, g, b uint32
}

// Struct for holding Least Significant Bit data
type Lsb struct {
    r, g, b uint8
}

// Struct for holding binary data in 8-bit chunks
type Bin struct {
    r, g, b uint8
}
