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

		defer func() {
	if err := os.Remove(tmpFile.Name()); err != nil {
		log.Println("failed to remove temp file:", err)
	}
}()


		if err := ttsFunc(text, tmpFile.Name()); err != nil {
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
	http.Error(w, "failed to stream file", http.StatusInternalServerError)
	return
}

	}
}


	func main() {
		http.HandleFunc("/speak", speakHandler(tts.TextToSpeech))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    if _, err := w.Write([]byte("OK")); err != nil {
    http.Error(w, "failed to write response", http.StatusInternalServerError)
    return
}

})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
