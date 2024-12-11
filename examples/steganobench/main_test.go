package main

import (
	"os"
	"testing"
)

func BenchmarkSteganoEmbed(b *testing.B) {
	b.StopTimer()
	data, err := os.ReadFile("data.txt")
	if err != nil {
		return
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := steganoEmbed("./in/big.png", data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAuyerEmbed(b *testing.B) {
	b.StopTimer()
	data, err := os.ReadFile("data.txt")
	if err != nil {
		return
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := auyerEmbed("./in/big.png", "out.png", data)
		if err != nil {
			b.Fatal(err)
		}
	}
}