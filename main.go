package main

import (
	"errors"
	"log"
	"os"

	"github.com/justSteven-lang/text-to-speech/tts"
)

// run contains business logic so it can be tested
func run(args []string) error {
	if len(args) < 2 {
		return errors.New("missing text argument")
	}

	text := args[1]
	return tts.TextToSpeech(text, "output.wav")
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}

	log.Println("Audio generated: output.wav")
}
