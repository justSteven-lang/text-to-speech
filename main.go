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

	tmpFile, err := os.CreateTemp("", "speech-*.wav")
	if err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}

	// Cleanup temp file safely
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			log.Println("failed to remove temp file:", err)
		}
	}()

	if err := tts.TextToSpeech(text, tmpFile.Name()); err != nil {
		http.Error(w, "failed to generate audio", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		http.Error(w, "failed to read audio file", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println("failed to close file:", err)
		}
	}()

	w.Header().Set("Content-Type", "audio/wav")

	if _, err := io.Copy(w, file); err != nil {
		log.Println("failed to write response:", err)
	}
}

func main() {
	http.HandleFunc("/speak", speakHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
