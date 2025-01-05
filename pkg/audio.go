package pkg

import (
	"github.com/go-audio/audio"
)


func EmbedDataAtDepthAudio(buffer *audio.IntBuffer, data []byte, depth uint8) *audio.IntBuffer {
	dataBits := BytesToBinary(data)
	lenBits := Int32ToBinary(int32(len(data)))
	lenBits = append(lenBits, dataBits...)

	for i := 0; i < len(lenBits); i++ {
		if buffer.Data[i] != int(lenBits[i]) {
			buffer.Data[i] = int(FlipBit(uint32(buffer.Data[i]), depth))
		}
	}

	return buffer
}

func ExtractDataAtDepthAudio(buffer *audio.IntBuffer, depth uint8) []byte {
	var data = make([]byte, 0)

	var currentByte uint8
	var count int
	for i := 0; i < len(buffer.Data); i++ {
		bit := GetBit(uint32(buffer.Data[i]), depth)
		currentByte = ((currentByte) << 1) | (bit & 1)
		count++

		if count == 8 {
			data = append(data, currentByte)
			currentByte = 0
			count = 0
		}
	}

	leng, _ := GetlenOfData(data)
	return data[4:(leng*2)+1]
}


func EmbedDataWithDepthAudio(buffer *audio.IntBuffer, data []byte, bitDepth uint8) *audio.IntBuffer {
	dataBits := BytesToBinary(data)
	lenBits := Int32ToBinary(int32(len(data)))
	lenBits = append(lenBits, dataBits...)

	curbit := bitDepth
	index := 0

	for i := 0; i < len(lenBits); i++ {
		if lenBits[i] != GetBit(uint32(buffer.Data[index]), curbit) {
			buffer.Data[index] = int(FlipBit(uint32(buffer.Data[index]), curbit))
		}
		
		if curbit != 0 {
			curbit--
		} else {
			curbit = bitDepth
			index++
		}
	}

	return buffer
}

func ExtractDataWithDepthAudio(buffer *audio.IntBuffer, depth uint8) []byte {
	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(buffer.Data); i++ {
		for bd := depth + 1; bd > 0; bd-- {
			r := GetBit(uint32(buffer.Data[i]), bd-1)
			currentByte = (currentByte << 1) | (r & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		currentByte = currentByte << (8 - bitCount)
		byteSlice = append(byteSlice, currentByte)
	}


	leng, _ := GetlenOfData(byteSlice)
	return byteSlice[4:(leng*2)+1]
}

