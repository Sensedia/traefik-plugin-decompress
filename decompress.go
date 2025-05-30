package traefik_plugin_decompress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

// Config holds the plugin configuration.
type Config struct{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// DecompressMiddleware is a plugin that decompresses gzip responses.
type DecompressMiddleware struct {
	next http.Handler
}

// New creates a new DecompressMiddleware instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &DecompressMiddleware{next: next}, nil
}

func (m *DecompressMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	slog.Info("DecompressMiddleware: starting request processing", "method", req.Method, "url", req.URL.String())

	if req.Header.Get("x-sensedia-gzip") == "true" {
		slog.Info("DecompressMiddleware: gzip encoding detected")
		contentType := req.Header.Get("x-sensedia-content-type")

		gr, err := gzip.NewReader(req.Body)
		if err != nil {
			slog.Error("DecompressMiddleware: failed to create gzip reader", "error", err)
			http.Error(rw, "Error - DecompressMiddleware: Failed to decompress request body", http.StatusBadRequest)
			return
		}
		defer gr.Close()

		var decompressed bytes.Buffer
		if _, err := io.Copy(&decompressed, gr); err != nil {
			slog.Error("DecompressMiddleware: failed to decompress body", "error", err)
			http.Error(rw, "Error - DecompressMiddleware: Failed to read decompressed data", http.StatusBadRequest)
			return
		}

		slog.Info("DecompressMiddleware: decompressed body size", "size", decompressed.Len())

		req.Body = io.NopCloser(bytes.NewReader(decompressed.Bytes()))
		req.ContentLength = int64(decompressed.Len())
		req.Header.Del("Content-Encoding")
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Content-Length", strconv.Itoa(decompressed.Len()))
	} else {
		slog.Info("DecompressMiddleware: no gzip encoding detected")
	}

	slog.Info("DecompressMiddleware: passing request to next handler")
	m.next.ServeHTTP(rw, req)
}
