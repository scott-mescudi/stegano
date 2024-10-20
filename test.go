package main

import "fmt"

func int32ToBinary(num int32) []int {
    var bits []int

    // Iterate over each bit position from 31 to 0
    for i := 31; i >= 0; i-- {
        bit := (num >> i) & 1 // Extract each bit from the integer
        bits = append(bits, int(bit)) // Append the bit to the slice
    }

    return bits
}

func packi32(n []int) []byte{
    var byteSlice = []byte{}
    var currentByte uint8 = 0
    bitCount := 0

    for _, i := range n {
        currentByte = (currentByte << 1) | (uint8(i) & 1)
        bitCount++

        if bitCount == 8 {
            byteSlice = append(byteSlice, currentByte)
            currentByte = 0
            bitCount = 0
        }
    }

    // Append any remaining bits in currentByte if bitCount is not zero
    if bitCount > 0 {
        // Shift currentByte to the left to fill the unused bits with zeros
        currentByte <<= (8 - bitCount)
        byteSlice = append(byteSlice, currentByte)
    }
	
	// Loop through the slice in increments of 8
	for i := 0; i < len(n); i += 8 {
		// Determine the end index for slicing
		end := i + 8
		if end > len(n) {
			end = len(n)
		}
		// Print the slice from i to end
		fmt.Println(n[i:end])
	}

	return byteSlice
}



func main() {
	var n int32 = 1_000_000

	binary := int32ToBinary(n)
	packed := packi32(binary)
	fmt.Printf("Binary: %v\n", binary)
	fmt.Printf("Packed:  %b\n", packed)

}