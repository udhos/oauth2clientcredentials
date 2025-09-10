package clientcredentials

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	const (
		clientID     = "myclientid"
		clientSecret = "myclientsecret"
		requestScope = "read write"
		accessToken  = "myaccesstoken"
		expiresIn    = 3600
	)

	// start http test server

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected content type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		req, err := DecodeRequestBody(r)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.GrantType != "client_credentials" {
			t.Errorf("expected grant_type client_credentials, got %s", req.GrantType)
		}
		if req.ClientID != clientID {
			t.Errorf("expected client_id %s, got %s", clientID, req.ClientID)
		}
		if req.ClientSecret != clientSecret {
			t.Errorf("expected client_secret %s, got %s", clientSecret, req.ClientSecret)
		}

		scope := req.Scope
		if scope != requestScope {
			t.Errorf("expected scope '%s', got %s", requestScope, scope)
		}

		// respond with access token
		respBody := EncodeResponseBody(accessToken, expiresIn)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respBody))
	}))

	resp, errSend := SendRequest(
		context.TODO(),
		http.DefaultClient,
		server.URL,
		clientID,
		clientSecret,
		requestScope,
	)
	if errSend != nil {
		t.Fatalf("failed to send request: %v", errSend)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		t.Fatalf("failed to read response body: %v", errBody)
	}

	server.Close()

	tokenResp, errDecode := DecodeResponseBody(body)
	if errDecode != nil {
		t.Fatalf("failed to decode response body: %v", errDecode)
	}

	if tokenResp.TokenType != "Bearer" {
		t.Errorf("expected token_type Bearer, got %s", tokenResp.TokenType)
	}

	if tokenResp.AccessToken != accessToken {
		t.Errorf("expected access_token %s, got %s", accessToken, tokenResp.AccessToken)
	}

	if tokenResp.ExpiresIn != expiresIn {
		t.Errorf("expected expires_in %d, got %d", expiresIn, tokenResp.ExpiresIn)
	}
}

type requestBodyTest struct {
	name         string
	clientID     string
	clientSecret string
	scope        string
	expected     string
}

var requestBodyTests = []requestBodyTest{
	{
		name:         "normal case",
		clientID:     "myclientid",
		clientSecret: "myclientsecret",
		scope:        "read write",
		expected:     "client_id=myclientid&client_secret=myclientsecret&grant_type=client_credentials&scope=read+write",
	},
	{
		name:         "unicode",
		clientID:     "myclientidá",
		clientSecret: "myclientsecretá",
		scope:        "á",
		expected:     "client_id=myclientid%C3%A1&client_secret=myclientsecret%C3%A1&grant_type=client_credentials&scope=%C3%A1",
	},
	{
		name:         "empty scope",
		clientID:     "myclientid",
		clientSecret: "myclientsecret",
		scope:        "",
		expected:     "client_id=myclientid&client_secret=myclientsecret&grant_type=client_credentials",
	},
}

// go test -count 1 -run ^TestEncodeRequestBody$ ./...
func TestEncodeRequestBody(t *testing.T) {
	for i, data := range requestBodyTests {
		name := fmt.Sprintf("%d of %d: %s", i+1, len(requestBodyTests), data.name)
		t.Run(name, func(t *testing.T) {
			result := EncodeRequestBody(data.clientID, data.clientSecret, data.scope)
			if result != data.expected {
				t.Fatalf("expected '%s', got '%s'", data.expected, result)
			}

			req, errReq := http.NewRequest("POST", "http://example.com/token", io.NopCloser(strings.NewReader(result)))
			if errReq != nil {
				t.Fatalf("failed to create request: %v", errReq)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			decodedReq, errDecode := DecodeRequestBody(req)
			if errDecode != nil {
				t.Fatalf("ailed to decode request body: %v", errDecode)
			}

			if decodedReq.GrantType != "client_credentials" {
				t.Errorf("expected grant_type client_credentials, got %s", decodedReq.GrantType)
			}
			if decodedReq.ClientID != data.clientID {
				t.Errorf("expected client_id %s, got %s", data.clientID, decodedReq.ClientID)
			}
			if decodedReq.ClientSecret != data.clientSecret {
				t.Errorf("expected client_secret %s, got %s", data.clientSecret, decodedReq.ClientSecret)
			}
			if decodedReq.Scope != data.scope {
				t.Errorf("expected scope %s, got %s", data.scope, decodedReq.Scope)
			}
		})
	}
}
