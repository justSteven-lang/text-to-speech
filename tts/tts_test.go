package tts

import (
	"os"
	"testing"
)

func TestTextToSpeech(t *testing.T) {
	output := "test.wav"

	err := TextToSpeech("hello test", output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatalf("expected %s to be created", output)
	}

	// cleanup
	_ = os.Remove(output)
}
