package game

import (
	"log"
	"os"
)

var AethelgardFont []byte

func init() {
	// Load font file directly
	data, err := os.ReadFile("assets/AethelgardFont.ttf")
	if err != nil {
		log.Fatal("Failed to load font file: ", err)
	}
	AethelgardFont = data
}
