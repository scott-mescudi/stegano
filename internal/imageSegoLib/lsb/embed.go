package lsb


// embedIntoRGBchannels embeds the binary data into the image's RGB channels using LSB technique.
func EmbedIntoRGBchannels(RGBchannels []rgbChannel, data []byte) []rgbChannel {
    z := splitIntoGroupsOfThree(bytesToBinary(data))
    
    for i := 0; i < len(z); i++ {
        if z[i].r != getLSB(RGBchannels[i].r) {
            RGBchannels[i].r = flipLSB(RGBchannels[i].r)
        }

        if z[i].g != getLSB(RGBchannels[i].g) {
            RGBchannels[i].g = flipLSB(RGBchannels[i].g)
        }

        if z[i].b != getLSB(RGBchannels[i].b) {
            RGBchannels[i].b = flipLSB(RGBchannels[i].b)
        }
    }

    return RGBchannels
}
