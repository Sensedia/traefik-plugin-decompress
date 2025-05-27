package traefik_plugin_decompress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Helper para criar um body gzip
func gzipCompress(data []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write(data)
	_ = gz.Close()
	return buf.Bytes()
}

func TestDecompressMiddleware_GzipRequest(t *testing.T) {
	originalBody := []byte("hello world")
	compressedBody := gzipCompress(originalBody)

	// Handler mock para verificar se o body foi descomprimido
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}
		if string(body) != string(originalBody) {
			t.Errorf("Expected body %q, got %q", originalBody, body)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware, _ := New(context.Background(), handler, CreateConfig(), "test")

	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(compressedBody))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Length", strconv.Itoa(len(compressedBody)))

	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestDecompressMiddleware_NonGzipRequest(t *testing.T) {
	originalBody := []byte("plain text")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}
		if string(body) != string(originalBody) {
			t.Errorf("Expected body %q, got %q", originalBody, body)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware, _ := New(context.Background(), handler, CreateConfig(), "test")

	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(originalBody))
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestDecompressMiddleware_InvalidGzip(t *testing.T) {
	invalidGzip := []byte("not really gzip")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called on invalid gzip")
	})

	middleware, _ := New(context.Background(), handler, CreateConfig(), "test")

	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(invalidGzip))
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %d", rr.Code)
	}
}
