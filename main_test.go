package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpeakHandler_MissingText(t *testing.T) {
	req := httptest.NewRequest("GET", "/speak", nil)
	w := httptest.NewRecorder()

	speakHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSpeakHandler_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/speak?text=hello", nil)
	w := httptest.NewRecorder()

	speakHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
