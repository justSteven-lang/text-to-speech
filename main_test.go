package main

import "testing"

func TestTextToSpeech(t *testing.T) {
	input := "hello devops"
	expected := "AUDIO_SIMULATION: hello devops"

	result, err := textToSpeech(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
