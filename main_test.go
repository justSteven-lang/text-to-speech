package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSpeakHandler_MissingText(t *testing.T) {
	mockTTS := func(text, filename string) error {
		return nil
	}

	handler := speakHandler(mockTTS)

	req := httptest.NewRequest(http.MethodGet, "/speak", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}


func TestSpeakHandler_Success(t *testing.T) {
	mockTTS := func(text, filename string) error {
		return os.WriteFile(filename, []byte("dummy audio"), 0644)
	}

	handler := speakHandler(mockTTS)

	req := httptest.NewRequest(http.MethodGet, "/speak?text=hello", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "audio/wav" {
		t.Fatalf("expected audio/wav content-type")
	}

	if rr.Body.Len() == 0 {
		t.Fatalf("expected body to contain audio data")
	}
}



func TestSpeakHandler_TTSFail(t *testing.T) {
	mockTTS := func(text, filename string) error {
		return errors.New("tts failed")
	}

	handler := speakHandler(mockTTS)

	req := httptest.NewRequest(http.MethodGet, "/speak?text=hello", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}


func TestHealthEndpoint(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if rr.Body.String() != "OK" {
		t.Fatalf("expected OK body")
	}
}


func TestNewServer_Health(t *testing.T) {
	server := newServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}


