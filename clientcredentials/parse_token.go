package clientcredentials

import (
	"fmt"
	"strconv"

	"github.com/sugawarayuuta/sonnet"
)

func parseToken(buf []byte, debugf func(format string, v ...any)) (Response, error) {
	var info Response

	var data map[string]interface{}

	errJSON := sonnet.Unmarshal(buf, &data)
	if errJSON != nil {
		return info, errJSON
	}

	accessToken, foundToken := data["access_token"]
	if !foundToken {
		return info, fmt.Errorf("missing access_token field in token response")
	}

	tokenStr, isStr := accessToken.(string)
	if !isStr {
		return info, fmt.Errorf("non-string value for access_token field in token response")
	}

	if tokenStr == "" {
		return info, fmt.Errorf("empty access_token in token response")
	}

	tokenType, foundTokenType := data["token_type"]
	if foundTokenType {
		tokenTypeStr, isStr := tokenType.(string)
		if !isStr {
			return info, fmt.Errorf("non-string value for token_type field in token response")
		}
		info.TokenType = tokenTypeStr
	}

	info.AccessToken = tokenStr

	expire, foundExpire := data["expires_in"]
	if foundExpire {
		switch expireVal := expire.(type) {
		case float64:
			debugf("found expires_in field with %f seconds", expireVal)
			info.ExpiresIn = int(expireVal)
		case string:
			debugf("found expires_in field with %s seconds", expireVal)
			exp, errConv := strconv.Atoi(expireVal)
			if errConv != nil {
				return info, fmt.Errorf("error converting expires_in field from string='%s' to int: %v", expireVal, errConv)
			}
			info.ExpiresIn = exp
		default:
			return info, fmt.Errorf("unexpected type %T for expires_in field in token response", expire)
		}
	}

	return info, nil
}
