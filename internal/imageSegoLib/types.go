package imagesegolib

// Struct for holding RGB channel data as 32-bit
type RgbChannel struct {
    R, G, B uint32
}

// Struct for holding Least Significant Bit data
type Lsb struct {
    R, G, B uint8
}

// Struct for holding binary data in 8-bit chunks
type Bin struct {
    R, G, B uint8
}
