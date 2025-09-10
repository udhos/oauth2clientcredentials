// Package clientcredentials provides functions to encode and decode
// client credentials token requests and responses.
package clientcredentials

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sugawarayuuta/sonnet"
)

// EncodeRequestBodyOld encodes the request body for client credentials grant type.
func EncodeRequestBodyOld(clientID, clientSecret, scope string) string {

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	if scope != "" {
		form.Add("scope", scope)
	}

	return form.Encode()
}

var (
	clientIDEncoded     = url.QueryEscape("client_id")
	clientSecretEncoded = "&" + url.QueryEscape("client_secret")
	grantTypeEncoded    = "&" + url.QueryEscape("grant_type") + "=" + url.QueryEscape("client_credentials")
	scopeEncoded        = "&" + url.QueryEscape("scope")
)

// EncodeRequestBody encodes the request body for client credentials grant type.
func EncodeRequestBody(clientID, clientSecret, scope string) string {

	if scope != "" {
		return clientIDEncoded + "=" + url.QueryEscape(clientID) + clientSecretEncoded + "=" + url.QueryEscape(clientSecret) + grantTypeEncoded + scopeEncoded + "=" + url.QueryEscape(scope)
	}

	return clientIDEncoded + "=" + url.QueryEscape(clientID) + clientSecretEncoded + "=" + url.QueryEscape(clientSecret) + grantTypeEncoded
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

// DecodeResponseBodyOld decodes the response body for client credentials grant type.
func DecodeResponseBodyOld(data []byte) (Response, error) {
	return parseToken(data, func(format string, v ...any) {})
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

// RequestOptions contains options for sending a client credentials token request.
type RequestOptions struct {
	// HTTPClient is optional HTTP client to use for sending the request.
	// If nil, http.DefaultClient will be used.
	HTTPClient HTTPDoer

	TokenURL     string
	ClientID     string
	ClientSecret string
	Scope        string

	// IsStatusCodeOK is optional function to check if the status code is OK.
	// If nil, DefaultIsStatusCodeOK will be used.
	IsStatusCodeOK func(statusCode int) bool
}

// DefaultIsStatusCodeOK is the default implementation for checking if a status code is OK.
func DefaultIsStatusCodeOK(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// SendRequest sends a client credentials token request and returns the response.
func SendRequest(ctx context.Context, options RequestOptions) (Response, error) {

	var tokenResp Response

	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}

	if options.IsStatusCodeOK == nil {
		options.IsStatusCodeOK = DefaultIsStatusCodeOK
	}

	reqBody := EncodeRequestBody(options.ClientID, options.ClientSecret, options.Scope)

	req, errReq := http.NewRequestWithContext(ctx, "POST", options.TokenURL, strings.NewReader(reqBody))
	if errReq != nil {
		return tokenResp, errReq
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, errDo := options.HTTPClient.Do(req)
	if errDo != nil {
		return tokenResp, errDo
	}

	defer resp.Body.Close()

	if !options.IsStatusCodeOK(resp.StatusCode) {
		return tokenResp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return tokenResp, errRead
	}

	return DecodeResponseBody(body)
}
