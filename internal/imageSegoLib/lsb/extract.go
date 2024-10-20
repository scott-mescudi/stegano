package lsb

type rgbChannel struct {
    r, g, b uint32
}


// extractDataFromRGBchannels extracts the hidden data from the image's RGB channels.
func ExtractDataFromRGBchannels(RGBchannels []rgbChannel) []byte {
    var byteSlice = make([]byte, 0)
    var currentByte uint8 = 0
    bitCount := 0

    for i := 0; i < len(RGBchannels); i++ {
        r := getLSB(RGBchannels[i].r)
        currentByte = (currentByte << 1) | (r & 1)
        bitCount++

        if bitCount == 8 {
            byteSlice = append(byteSlice, currentByte)
            currentByte = 0
            bitCount = 0
        }

        g := getLSB(RGBchannels[i].g)
        currentByte = (currentByte << 1) | (g & 1)
        bitCount++

        if bitCount == 8 {
            byteSlice = append(byteSlice, currentByte)
            currentByte = 0
            bitCount = 0
        }

        b := getLSB(RGBchannels[i].b)
        currentByte = (currentByte << 1) | (b & 1)
        bitCount++

        if bitCount == 8 {
            byteSlice = append(byteSlice, currentByte)
            currentByte = 0
            bitCount = 0
        }
    }

    if bitCount > 0 {
        currentByte = currentByte << (8 - bitCount)
        byteSlice = append(byteSlice, currentByte)
    }

    return byteSlice
}
