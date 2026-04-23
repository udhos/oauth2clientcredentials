[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/oauth2clientcredentials/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/oauth2clientcredentials)](https://goreportcard.com/report/github.com/udhos/oauth2clientcredentials)
[![Go Reference](https://pkg.go.dev/badge/github.com/udhos/oauth2clientcredentials.svg)](https://pkg.go.dev/github.com/udhos/oauth2clientcredentials)

# oauth2clientcredentials

Package clientcredentials provides functions to encode and decode oauth2 client credentials token requests and responses.

# Synopsis

## Client

See full client example: [examples/clientcredentials-token-client/main.go](examples/clientcredentials-token-client/main.go)

```go
import "github.com/udhos/oauth2clientcredentials/clientcredentials"

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
```

## Server

See full server example: [examples/clientcredentials-token-server/main.go](examples/clientcredentials-token-server/main.go)

```go
import "github.com/udhos/oauth2clientcredentials/clientcredentials"

func handlerToken(w http.ResponseWriter, r *http.Request) {

	req, err := clientcredentials.DecodeRequestBody(r)

    replyStr = clientcredentials.EncodeResponseBody(accessToken, scope,	expireSeconds)
```

# References

- [RFC6749 The OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
