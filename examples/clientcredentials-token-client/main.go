// Package main implements the tool.
package main

import (
	"context"
	"flag"
	"log"

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

	options := clientcredentials.RequestOptions{
		TokenURL:     tokenURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	tokenResp, errSend := clientcredentials.SendRequest(context.TODO(), options)
	if errSend != nil {
		log.Fatalf("error http post: %v", errSend)
	}

	// print response fields
	log.Printf("access_token: %s", tokenResp.AccessToken)
	log.Printf("token_type: %s", tokenResp.TokenType)
	log.Printf("expires_in: %d", tokenResp.ExpiresIn)
	log.Printf("scope: %s", tokenResp.Scope)
}
