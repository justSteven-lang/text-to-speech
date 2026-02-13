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
		// create dummy file so open doesn't fail
		return os.WriteFile(filename, []byte("dummy audio"), 0644)
	}

	handler := speakHandler(mockTTS)

	req := httptest.NewRequest(http.MethodGet, "/speak?text=hello", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
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


