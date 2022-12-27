package tokens

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"

	"github.com/Nigel2392/jsext"
)

func SetTokenCookie(token *Token) {
	var b bytes.Buffer
	// Gob the token
	var enc = gob.NewEncoder(&b)
	enc.Encode(token)
	// Encode to base64
	var cookie = base64.StdEncoding.EncodeToString(b.Bytes())
	// Set the cookie
	jsext.SetCookie("token", cookie, 1*3600*24)
}

func GetTokenCookie() (*Token, error) {
	// Get the cookie
	var cookie = jsext.GetCookie("token")
	if cookie == "" {
		return nil, errors.New("no token cookie")
	}
	// Decode from base64
	var b, err = base64.StdEncoding.DecodeString(cookie)
	if err != nil {
		return nil, err
	}
	// Gob the token
	var dec = gob.NewDecoder(bytes.NewReader(b))
	var token Token
	dec.Decode(&token)
	token.updateManager()
	return &token, nil
}
