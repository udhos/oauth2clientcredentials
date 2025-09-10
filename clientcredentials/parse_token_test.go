package clientcredentials

import (
	"testing"
)

const (
	expectSucess  = true
	expectFailure = false
)

type expectResult bool

type parseTokenTestCase struct {
	name             string
	token            string
	expect           expectResult
	expectAcessToken string
	expectTokenType  string
	expectExpire     int
}

var parseTokenTestTable = []parseTokenTestCase{
	{"empty", "", expectFailure, "", "", 0},
	{"not-json", "not-json", expectFailure, "", "", 0},
	{"no fields", `{}`, expectFailure, "", "", 0},
	{"missing access_token", `{"other":"field"}`, expectFailure, "", "", 0},
	{"empty access_token", `{"access_token":""}`, expectFailure, "", "", 0},
	{"only good access token", `{"access_token":"abc"}`, expectSucess, "abc", "", 0},
	{"token type", `{"access_token":"abc","token_type":"Bearer"}`, expectSucess, "abc", "Bearer", 0},
	{"wrong access token type int", `{"access_token":123}`, expectFailure, "", "", 0},
	{"wrong access token type bool", `{"access_token":true}`, expectFailure, "", "", 0},
	{"wrong access token type float", `{"access_token":1.1}`, expectFailure, "", "", 0},
	{"expire integer", `{"access_token":"abc","expires_in":300}`, expectSucess, "abc", "", 300},
	{"expire float", `{"access_token":"abc","expires_in":300.0}`, expectSucess, "abc", "", 300},
	{"expire string", `{"access_token":"abc","expires_in":"300"}`, expectSucess, "abc", "", 300},
	{"expire broken string", `{"access_token":"abc","expires_in":"TTT"}`, expectFailure, "", "", 0},
	{"expire empty string", `{"access_token":"abc","expires_in":""}`, expectFailure, "", "", 0},
	{"expire broken bool", `{"access_token":"abc","expires_in":true}`, expectFailure, "", "", 0},
}

func TestParseToken(t *testing.T) {
	for _, data := range parseTokenTestTable {
		buf := []byte(data.token)
		info, errParse := parseToken(buf, t.Logf)
		success := errParse == nil
		if success != bool(data.expect) {
			t.Errorf("%s: expectedSuccess=%t gotSuccess=%t error:%v", data.name, data.expect, success, errParse)
			continue
		}

		if !success {
			continue
		}

		var errored bool

		if info.AccessToken != data.expectAcessToken {
			t.Errorf("%s: expectedAccessToken=%s gotAccessToken=%s", data.name, data.expectAcessToken, info.AccessToken)
			errored = true
		}

		if info.ExpiresIn != data.expectExpire {
			t.Errorf("%s: expectedExpire=%v gotExpire=%v", data.name, data.expectExpire, info.ExpiresIn)
			errored = true
		}

		if info.TokenType != data.expectTokenType {
			t.Errorf("%s: expectedTokenType=%s gotTokenType=%s", data.name, data.expectTokenType, info.TokenType)
			errored = true
		}

		if !errored {
			t.Logf("%s: ok", data.name)
		}
	}
}
