package main

import (
	"bytes"
	"container/heap"
	"encoding/gob"
	"fmt"
	"strings"
)

// HuffmanNode represents a node in the Huffman tree
type HuffmanNode struct {
	Char  rune         // Character (for leaf nodes)
	Freq  int          // Frequency of the character
	Left  *HuffmanNode // Left child
	Right *HuffmanNode // Right child
}

// PriorityQueue implements heap.Interface for the Huffman Tree nodes
type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*HuffmanNode))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// BuildFrequencyTable counts the frequency of each character in the input string
func BuildFrequencyTable(input string) map[rune]int {
	frequency := make(map[rune]int)
	for _, char := range input {
		frequency[char]++
	}
	return frequency
}

// BuildHuffmanTree constructs the Huffman Tree using the frequency table
func BuildHuffmanTree(frequency map[rune]int) *HuffmanNode {
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Create a leaf node for each character and add it to the priority queue
	for char, freq := range frequency {
		heap.Push(pq, &HuffmanNode{Char: char, Freq: freq})
	}

	// Merge nodes until there's only one node (the root of the Huffman tree)
	for pq.Len() > 1 {
		left := heap.Pop(pq).(*HuffmanNode)
		right := heap.Pop(pq).(*HuffmanNode)

		// Create a new internal node with frequency equal to the sum of two nodes
		newNode := &HuffmanNode{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(pq, newNode)
	}

	// The last remaining node is the root of the Huffman tree
	return heap.Pop(pq).(*HuffmanNode)
}

// GenerateCodes generates a Huffman code for each character by traversing the tree
func GenerateCodes(node *HuffmanNode, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = prefix
	}
	GenerateCodes(node.Left, prefix+"0", codes)
	GenerateCodes(node.Right, prefix+"1", codes)
}

// Encode converts the input string into its Huffman encoded binary string
func Encode(input string, codes map[rune]string) string {
	var encoded strings.Builder
	for _, char := range input {
		encoded.WriteString(codes[char])
	}
	return encoded.String()
}

// Decode decodes a Huffman encoded binary string back into the original string
func Decode(encoded string, root *HuffmanNode) string {
	var decoded strings.Builder
	current := root

	for _, bit := range encoded {
		if bit == '0' {
			current = current.Left
		} else {
			current = current.Right
		}

		// If it's a leaf node, append the character to the result
		if current.Left == nil && current.Right == nil {
			decoded.WriteRune(current.Char)
			current = root
		}
	}

	return decoded.String()
}

// Serialize the Huffman tree using GOB (Go Binary Encoding)
func SerializeTree(root *HuffmanNode) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(root)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize the Huffman tree from the binary format
func DeserializeTree(data []byte) (*HuffmanNode, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var root HuffmanNode
	err := decoder.Decode(&root)
	if err != nil {
		return nil, err
	}
	return &root, nil
}

// Pack binary string into byte array
func PackBinaryString(binaryStr string) []byte {
	var buf bytes.Buffer
	byteVal := byte(0)
	bitCount := 0

	for _, bit := range binaryStr {
		byteVal <<= 1
		if bit == '1' {
			byteVal |= 1
		}
		bitCount++
		if bitCount == 8 {
			buf.WriteByte(byteVal)
			byteVal = 0
			bitCount = 0
		}
	}

	// If there's remaining bits, pad the rest
	if bitCount > 0 {
		byteVal <<= (8 - bitCount)
		buf.WriteByte(byteVal)
	}

	return buf.Bytes()
}

// Unpack byte array back into binary string
func UnpackBinaryString(encoded []byte) string {
	var binaryStr strings.Builder
	for _, b := range encoded {
		for i := 7; i >= 0; i-- {
			if (b>>i)&1 == 1 {
				binaryStr.WriteByte('1')
			} else {
				binaryStr.WriteByte('0')
			}
		}
	}
	return binaryStr.String()
}

// HuffmanEncode takes input text and returns serialized tree data and packed data
func HuffmanEncode(input string) ([]byte, []byte, error) {
	frequencyTable := BuildFrequencyTable(input)
	huffmanTree := BuildHuffmanTree(frequencyTable)

	// Step 3: Generate Huffman Codes
	codes := make(map[rune]string)
	GenerateCodes(huffmanTree, "", codes)

	// Step 4: Encode the input string
	encoded := Encode(input, codes)
	// Step 5: Pack encoded binary string into bytes for transfer
	packedData := PackBinaryString(encoded)

	// Step 6: Serialize the Huffman tree
	treeData, err := SerializeTree(huffmanTree)
	if err != nil {
		return nil, nil, fmt.Errorf("error serializing tree: %v", err)
	}

	return treeData, packedData, nil
}

// HuffmanDecode takes serialized tree data and packed data and returns the decoded string
func HuffmanDecode(treeData []byte, packedData []byte) (string, error) {
	receivedTree, err := DeserializeTree(treeData)
	if err != nil {
		return "", fmt.Errorf("error deserializing tree: %v", err)
	}

	// Step 10: Unpack the binary string from the received byte array
	unpackedBinary := UnpackBinaryString(packedData)

	decoded := Decode(unpackedBinary, receivedTree)
	return decoded, nil
}

// func main() {
// 	// Example usage:
// 	input := "this is an example for huffman encoding"
// 	treeData, packedData, err := HuffmanEncode(input)
// 	if err != nil {
// 		fmt.Println("Encoding error:", err)
// 		return
// 	}

// 	fmt.Println(len(treeData) + len(packedData))

// 	decoded, err := HuffmanDecode(treeData, packedData)
// 	if err != nil {
// 		fmt.Println("Decoding error:", err)
// 		return
// 	}

// 	fmt.Println("Original:", input)
// 	fmt.Println("Decoded:", decoded)
// }


