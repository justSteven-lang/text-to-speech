package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/justSteven-lang/text-to-speech/tts"
)

func speakHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "missing text parameter", http.StatusBadRequest)
		return
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "speech-*.wav")
	if err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Generate speech into temp file
	err = tts.TextToSpeech(text, tmpFile.Name())
	if err != nil {
		http.Error(w, "failed to generate audio", http.StatusInternalServerError)
		return
	}

	// Open file again to send response
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		http.Error(w, "failed to read audio file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/wav")
	w.WriteHeader(http.StatusOK)

	io.Copy(w, file)
}

func main() {
	http.HandleFunc("/speak", speakHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
