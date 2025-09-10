// Package clientcredentials provides functions to encode and decode
// client credentials token requests and responses.
package clientcredentials

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sugawarayuuta/sonnet"
)

// EncodeRequestBody encodes the request body for client credentials grant type.
func EncodeRequestBody(clientID, clientSecret, scope string) string {

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	if scope != "" {
		form.Add("scope", scope)
	}

	return form.Encode()
}

// DecodeRequestBody decodes the request body for client credentials grant type.
func DecodeRequestBody(r *http.Request) (Request, error) {

	var req Request

	if err := r.ParseForm(); err != nil {
		return req, err
	}

	req.GrantType = getParam(r, "grant_type")
	req.ClientID = getParam(r, "client_id")
	req.ClientSecret = getParam(r, "client_secret")
	req.Scope = getParam(r, "scope")

	return req, nil
}

// Request represents a client credentials token request.
type Request struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	Scope        string
}

func getParam(r *http.Request, key string) string {
	v := r.Form[key]
	if v == nil {
		return ""
	}
	return v[0]
}

// EncodeResponseBody encodes the response body for client credentials grant type.
func EncodeResponseBody(accessToken string, expiresInSeconds int) string {
	expiresInSecondsStr := strconv.Itoa(expiresInSeconds)
	return `{"access_token":"` + accessToken + `","token_type":"Bearer","expires_in":` + expiresInSecondsStr + `}`
}

// DecodeResponseBody decodes the response body for client credentials grant type.
func DecodeResponseBody(data []byte) (Response, error) {
	var resp Response
	err := sonnet.Unmarshal(data, &resp)
	return resp, err
}

// Response represents a client credentials token response.
type Response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

// HTTPDoer is an interface for plugging in custom HTTP clients.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendRequest sends a client credentials token request and returns the response.
func SendRequest(ctx context.Context, httpClient HTTPDoer, tokenURL,
	clientID, clientSecret, scope string) (*http.Response, error) {

	reqBody := EncodeRequestBody(clientID, clientSecret, scope)

	req, errReq := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(reqBody))
	if errReq != nil {
		log.Fatalf("error new request: %v", errReq)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return httpClient.Do(req)
}
