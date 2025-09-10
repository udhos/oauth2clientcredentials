// Package bench provides benchmarking utilities.
package bench

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/udhos/oauth2clientcredentials/clientcredentials"
)

// go test -bench=. ./bench
func BenchmarkEncodeRequestBody(b *testing.B) {
	for b.Loop() {
		clientcredentials.EncodeRequestBody("myclientid", "myclientsecret", "read write")
	}
}

// go test -bench=. ./bench
func BenchmarkDecodeRequestBody(b *testing.B) {
	reqBody := clientcredentials.EncodeRequestBody("myclientid", "myclientsecret", "read write")
	req, err := http.NewRequest("POST", "http://example.com/token", io.NopCloser(strings.NewReader(reqBody)))
	if err != nil {
		b.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for b.Loop() {
		_, err := clientcredentials.DecodeRequestBody(req)
		if err != nil {
			b.Fatalf("failed to decode request body: %v", err)
		}
	}
}

// go test -bench=. ./bench
func BenchmarkEncodeResponseBody(b *testing.B) {
	for b.Loop() {
		clientcredentials.EncodeResponseBody("myaccesstoken", 3600)
	}
}

// go test -bench=. ./bench
func BenchmarkDecodeResponseBody(b *testing.B) {
	respBody := clientcredentials.EncodeResponseBody("myaccesstoken", 3600)
	data := []byte(respBody)

	for b.Loop() {
		_, err := clientcredentials.DecodeResponseBody(data)
		if err != nil {
			b.Fatalf("failed to decode response body: %v", err)
		}
	}
}
