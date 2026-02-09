package main

import (
	"fmt"
	"os"
)

func textToSpeech(text string) error {
	file, err := os.Create("output.txt")
	if err != nil {
		return  err
	}
	defer file.Close()

	_, err = file.WriteString("AUDIO_SIMULATION: " + text)
	return err
}

func main()  {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"your text\"")
		return
	}

	text := os.Args[1]

	err := textToSpeech(text)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Text converted to audio (simulated). File: output.txt")
}