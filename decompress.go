package traefik_plugin_decompress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log"
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
	log.Println("DecompressMiddleware: received request")

	if req.Header.Get("x-sensedia-gzip") == "true" {
		log.Println("DecompressMiddleware: gzip encoding detected")

		gr, err := gzip.NewReader(req.Body)
		if err != nil {
			log.Printf("DecompressMiddleware: failed to create gzip reader: %v", err)
			http.Error(rw, "Failed to decompress request body", http.StatusBadRequest)
			return
		}
		defer gr.Close()

		var decompressed bytes.Buffer
		if _, err := io.Copy(&decompressed, gr); err != nil {
			log.Printf("DecompressMiddleware: failed to decompress body: %v", err)
			http.Error(rw, "Failed to read decompressed data", http.StatusBadRequest)
			return
		}

		log.Printf("DecompressMiddleware: decompressed body size: %d bytes", decompressed.Len())

		req.Body = io.NopCloser(bytes.NewReader(decompressed.Bytes()))
		req.ContentLength = int64(decompressed.Len())
		req.Header.Del("Content-Encoding")
		req.Header.Set("Content-Length", strconv.Itoa(decompressed.Len()))
	} else {
		log.Println("DecompressMiddleware: no gzip encoding detected")
	}

	log.Println("DecompressMiddleware: passing request to next handler")
	m.next.ServeHTTP(rw, req)
}
