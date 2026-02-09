package tts

import (
	"fmt"
	"os/exec"
)

func TextToSpeech(text string, output string) error {
	if text == "" {
		return fmt.Errorf("text is empty")
	}

	cmd := exec.Command(
		"espeak",
		text,
		"-w",
		output,
	)

	return cmd.Run()
}
