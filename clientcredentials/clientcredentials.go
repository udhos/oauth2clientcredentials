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
	"sync"

	"github.com/sugawarayuuta/sonnet"
	"github.com/valyala/fastjson"
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
func EncodeResponseBody(accessToken, scope string, expiresInSeconds int) string {
	expiresInSecondsStr := strconv.Itoa(expiresInSeconds)
	return `{"access_token":"` + accessToken + `","token_type":"Bearer","expires_in":` + expiresInSecondsStr + `,"scope":"` + scope + `"}`
}

// DecodeResponseBody decodes the response body for client credentials
// grant type using the fastest decoder (currently fastjson with sync.Pool).
//
// 2026-04-22
//
// $ go test -bench=DecodeResponse -benchmem ./bench
// goos: linux
// goarch: amd64
// pkg: github.com/udhos/oauth2clientcredentials/bench
// cpu: 13th Gen Intel(R) Core(TM) i7-1360P
// BenchmarkDecodeResponseBody-16                    	 4552692	       256.3 ns/op	      32 B/op	       3 allocs/op
// BenchmarkDecodeResponseBodySonnet-16              	 3078936	       379.1 ns/op	     216 B/op	       6 allocs/op
// BenchmarkDecodeResponseBodyCustomParser-16        	 1000000	      1084 ns/op	     696 B/op	      22 allocs/op
// BenchmarkDecodeResponseBodyFastJSON-16            	 1000000	      1186 ns/op	    1560 B/op	      11 allocs/op
// BenchmarkDecodeResponseBodyFastJSONSyncPool-16    	 4591796	       252.6 ns/op	      32 B/op	       3 allocs/op
// PASS
// ok  	github.com/udhos/oauth2clientcredentials/bench	5.772s
func DecodeResponseBody(data []byte) (Response, error) {
	return DecodeResponseBodyFastJSONSyncPool(data)
}

// DecodeResponseBodyFastJSON decodes the response body using valyala/fastjson.
func DecodeResponseBodyFastJSON(data []byte) (Response, error) {
	var resp Response

	var p fastjson.Parser
	v, err := p.ParseBytes(data)
	if err != nil {
		return resp, err
	}

	// Extract fields
	resp.AccessToken = string(v.GetStringBytes("access_token"))
	resp.TokenType = string(v.GetStringBytes("token_type"))

	// GetInt returns 0 if the key is missing or not an integer
	resp.ExpiresIn = v.GetInt("expires_in")
	resp.Scope = string(v.GetStringBytes("scope"))

	return resp, nil
}

var parserPool = sync.Pool{
	New: func() any { return new(fastjson.Parser) },
}

// DecodeResponseBodyFastJSONSyncPool decodes the response body using valyala/fastjson with sync.Pool for parser reuse.
func DecodeResponseBodyFastJSONSyncPool(data []byte) (Response, error) {
	var resp Response

	// Get a parser from the pool
	p := parserPool.Get().(*fastjson.Parser)
	defer parserPool.Put(p)

	v, err := p.ParseBytes(data)
	if err != nil {
		return resp, err
	}

	// Direct assignment avoids reflection
	resp.AccessToken = string(v.GetStringBytes("access_token"))
	resp.TokenType = string(v.GetStringBytes("token_type"))
	resp.ExpiresIn = v.GetInt("expires_in")
	resp.Scope = string(v.GetStringBytes("scope"))

	return resp, nil
}

// DecodeResponseBodySonnet decodes the response body for client credentials grant type using sonnet.
func DecodeResponseBodySonnet(data []byte) (Response, error) {
	var resp Response
	err := sonnet.Unmarshal(data, &resp)
	return resp, err
}

// DecodeResponseBodyCustomParser decodes the response body for client credentials grant type using a custom parser.
func DecodeResponseBodyCustomParser(data []byte) (Response, error) {
	return parseToken(data, func(_ string, _ ...any) {})
}

// Response represents a client credentials token response.
type Response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
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
	IsStatusCodeOK func(statusCode int) error
}

// DefaultIsStatusCodeOK is the default implementation for checking if a status code is OK.
func DefaultIsStatusCodeOK(statusCode int) error {
	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("oauth2clientcredentials.DefaultIsStatusCodeOK: status code out of range 200-299: %d", statusCode)
	}
	return nil // ok
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

	req, errReq := http.NewRequestWithContext(ctx, "POST", options.TokenURL,
		strings.NewReader(reqBody))
	if errReq != nil {
		return tokenResp, errReq
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, errDo := options.HTTPClient.Do(req)
	if errDo != nil {
		return tokenResp, errDo
	}

	defer resp.Body.Close()

	if err := options.IsStatusCodeOK(resp.StatusCode); err != nil {
		return tokenResp, fmt.Errorf("oauth2clientcredentials.SendRequest error: %w", err)
	}

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return tokenResp, errRead
	}

	return DecodeResponseBody(body)
}
