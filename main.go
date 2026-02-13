package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/justSteven-lang/text-to-speech/tts"
)

type TTSFunc func(text, filename string) error

func speakHandler(ttsFunc TTSFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		defer os.Remove(tmpFile.Name())

		if err := ttsFunc(text, tmpFile.Name()); err != nil {
			http.Error(w, "failed to generate audio", http.StatusInternalServerError)
			return
		}

		file, err := os.Open(tmpFile.Name())
		if err != nil {
			http.Error(w, "failed to read audio file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "audio/wav")
		io.Copy(w, file)
	}
}


	func main() {
		http.HandleFunc("/speak", speakHandler(tts.TextToSpeech))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
