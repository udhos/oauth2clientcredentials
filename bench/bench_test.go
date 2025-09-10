// Package bench provides benchmarking utilities.
package bench

/*
Date: 2025-09-10

go version go1.25.1 linux/amd64

$ go test -bench=. ./bench
goos: linux
goarch: amd64
pkg: github.com/udhos/oauth2clientcredentials/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkEncodeRequestBodyOld-16     	 1373432	       875.7 ns/op
BenchmarkEncodeRequestBody-16        	 5795644	       184.7 ns/op
BenchmarkDecodeRequestBody-16        	25285114	        44.02 ns/op
BenchmarkEncodeResponseBody-16       	16006507	        78.54 ns/op
BenchmarkDecodeResponseBody-16       	 4261275	       271.3 ns/op
BenchmarkDecodeResponseBodyOld-16    	 1367126	       901.3 ns/op
PASS
ok  	github.com/udhos/oauth2clientcredentials/bench	7.041s
*/

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/udhos/oauth2clientcredentials/clientcredentials"
)

// go test -bench=. ./bench
func BenchmarkEncodeRequestBodyOld(b *testing.B) {
	for b.Loop() {
		clientcredentials.EncodeRequestBodyOld("myclientid", "myclientsecret", "read write")
	}
}

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

// go test -bench=. ./bench
func BenchmarkDecodeResponseBodyOld(b *testing.B) {
	respBody := clientcredentials.EncodeResponseBody("myaccesstoken", 3600)
	data := []byte(respBody)

	for b.Loop() {
		_, err := clientcredentials.DecodeResponseBodyOld(data)
		if err != nil {
			b.Fatalf("failed to decode response body: %v", err)
		}
	}
}
