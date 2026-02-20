package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
	handler := newMux()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestNewServer_Config(t *testing.T) {
	server := newServer()

	if server.Addr != ":8080" {
		t.Fatalf("expected addr :8080, got %s", server.Addr)
	}

	if server.Handler == nil {
		t.Fatalf("expected handler to be set")
	}

	if server.ReadTimeout == 0 {
		t.Fatalf("expected ReadTimeout to be set")
	}

	if server.WriteTimeout == 0 {
		t.Fatalf("expected WriteTimeout to be set")
	}

	if server.IdleTimeout == 0 {
		t.Fatalf("expected IdleTimeout to be set")
	}
}


func TestServer_StartAndShutdown(t *testing.T) {
	server := newServer()
	server.Addr = ":0"

	done := make(chan struct{})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("ListenAndServe error: %v", err)
		}
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown error: %v", err)
	}

	<-done
}



func TestRun_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- run(ctx)
	}()

	// beri waktu server start
	time.Sleep(200 * time.Millisecond)

	// trigger shutdown
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("run did not shutdown in time")
	}
}


func TestNewMux_SpeakRouteExists(t *testing.T) {
	handler := newMux()

	req := httptest.NewRequest(http.MethodGet, "/speak?text=test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", rr.Code)
	}
}









