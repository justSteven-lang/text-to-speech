package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justSteven-lang/text-to-speech/tts"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TTSFunc func(text, filename string) error

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tts_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tts_http_request_duration_seconds",
			Help:    "HTTP request duration.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
	})
}




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

func newMux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/speak", speakHandler(tts.TextToSpeech))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	return metricsMiddleware(mux)
}


func newServer() *http.Server {
	return &http.Server{
		Addr:         ":8080",
		Handler:      newMux(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}


func run(ctx context.Context) error {
	server := newServer()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}


func main() {
	log.Println("Server running on :8080")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
}


