package main

import (
	"log"
	"os"

	"github.com/justSteven-lang/text-to-speech/tts"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go \"your text here\"")
	}

	text := os.Args[1]

	err := tts.TextToSpeech(text, "output.wav")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Audio generated: output.wav")
}
