// Package bench provides benchmarking utilities.
package bench

/*
Date: 2025-12-12

go version go1.25.1 linux/amd64

go test -bench=. ./bench
goos: linux
goarch: amd64
pkg: github.com/udhos/oauth2clientcredentials/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkEncodeRequestBodyOld-16     	 1218193	       984.0 ns/op
BenchmarkEncodeRequestBody-16        	 6438307	       175.7 ns/op
BenchmarkDecodeRequestBody-16        	25867268	        42.59 ns/op
BenchmarkEncodeResponseBody-16       	12504793	        96.34 ns/op
BenchmarkDecodeResponseBody-16       	 3598549	       339.2 ns/op
BenchmarkDecodeResponseBodyOld-16    	 1000000	      1082 ns/op
PASS
ok  	github.com/udhos/oauth2clientcredentials/bench	6.948s
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
		clientcredentials.EncodeResponseBody("myaccesstoken", "scope", 3600)
	}
}

// go test -bench=. ./bench
func BenchmarkDecodeResponseBody(b *testing.B) {
	respBody := clientcredentials.EncodeResponseBody("myaccesstoken", "scope", 3600)
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
	respBody := clientcredentials.EncodeResponseBody("myaccesstoken", "scope", 3600)
	data := []byte(respBody)

	for b.Loop() {
		_, err := clientcredentials.DecodeResponseBodyOld(data)
		if err != nil {
			b.Fatalf("failed to decode response body: %v", err)
		}
	}
}
