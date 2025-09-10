// Package main implements the tool.
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/udhos/oauth2clientcredentials/clientcredentials"
)

func main() {

	var clientID string
	var clientSecret string
	var scope string
	var tokenURL string

	flag.StringVar(&clientID, "client_id", "admin", "client id")
	flag.StringVar(&clientSecret, "client_secret", "admin", "client secret")
	flag.StringVar(&scope, "scope", "scope1", "scope")
	flag.StringVar(&tokenURL, "token_url", "http://localhost:8080/token", "token url")
	flag.Parse()

	resp, errDo := clientcredentials.SendRequest(context.TODO(), http.DefaultClient, tokenURL, clientID, clientSecret, scope)
	if errDo != nil {
		log.Fatalf("error http post: %v", errDo)
	}
	defer resp.Body.Close()

	log.Printf("response status: %d", resp.StatusCode)

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Fatalf("error read body: %v", errBody)
	}

	tokenResp, errDecode := clientcredentials.DecodeResponseBody(body)
	if errDecode != nil {
		log.Fatalf("error decode response body: %v", errDecode)
	}

	// print response fields
	log.Printf("access_token: %s", tokenResp.AccessToken)
	log.Printf("token_type: %s", tokenResp.TokenType)
	log.Printf("expires_in: %d", tokenResp.ExpiresIn)
}
