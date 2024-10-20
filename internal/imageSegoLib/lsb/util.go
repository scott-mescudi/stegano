package lsb


// Struct for holding binary data in 8-bit chunks
type bin struct {
    r, g, b uint8
}

// getLSB returns the least significant bit of a 32-bit value.
func getLSB(value uint32) uint8 {
    return uint8(value & 1)
}

// flipLSB flips the least significant bit of a 32-bit number.
func flipLSB(num uint32) uint32 {
    return num ^ 1
}

// bytesToBinary converts a byte slice into a slice of bits.
func bytesToBinary(data []byte) []int {
    var bits []int

    for _, b := range data {
        for i := 7; i >= 0; i-- {
            bit := (b >> i) & 1
            bits = append(bits, int(bit))
        }
    }

    return bits
}

// splitIntoGroupsOfThree splits a slice of bits into groups of three.
func splitIntoGroupsOfThree(nums []int) []bin {
    var result []bin
    for i := 0; i < len(nums); i += 3 {
        var b bin
        if i < len(nums) {
            b.r = uint8(nums[i])
        }
        if i+1 < len(nums) {
            b.g = uint8(nums[i+1])
        }
        if i+2 < len(nums) {
            b.b = uint8(nums[i+2])
        }
        result = append(result, b)
    }
    return result
}
