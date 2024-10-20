package lsb


import (
	s "lsb/internal/imageSegoLib"
)

// embedIntoRGBchannels embeds the binary data into the image's RGB channels using LSB technique.
func EmbedIntoRGBchannels(RGBchannels []s.RgbChannel, data []byte) []s.RgbChannel {
    z := splitIntoGroupsOfThree(bytesToBinary(data))
    
    for i := 0; i < len(z); i++ {
        if z[i].R != getLSB(RGBchannels[i].R) {
            RGBchannels[i].R = flipLSB(RGBchannels[i].R)
        }

        if z[i].G != getLSB(RGBchannels[i].G) {
            RGBchannels[i].G = flipLSB(RGBchannels[i].G)
        }

        if z[i].B != getLSB(RGBchannels[i].B) {
            RGBchannels[i].B = flipLSB(RGBchannels[i].B)
        }
    }

    return RGBchannels
}
