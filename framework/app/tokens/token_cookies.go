//go:build js && wasm
// +build js,wasm

package tokens

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/Nigel2392/jsext"
)

const JSEXT_token = "JSEXT_token"

//go:generate msgp -tests=false
type SaveToken struct {
	AccessToken  string    `msg:"access_token"`
	RefreshToken string    `msg:"refresh_token"`
	LastUpdate   time.Time `msg:"last_update"`
}

// Set the token cookie from a token.
func SetTokenCookie(token *Token) error {
	var AccessToken = token.AccessToken
	var RefreshToken = token.RefreshToken
	var LastUpdate = token.LastUpdate

	var saveToken = SaveToken{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		LastUpdate:   LastUpdate,
	}

	var msgBytes, err = saveToken.MarshalMsg(nil)
	if err != nil {
		return err
	}

	// Encode to base64
	var cookie = base64.RawURLEncoding.EncodeToString(msgBytes)
	// Set the cookie
	return jsext.SetCookie(JSEXT_token, cookie, token.RefreshExpiredIn())
}

// Get the token cookie
func GetTokenCookie(tokenToSet *Token) (*Token, error) {
	var cookie = jsext.GetCookie(JSEXT_token)
	if cookie == "" {
		//lint:ignore ST1005 Error strings should not be capitalized
		return nil, errors.New("No token cookie")
	}

	var saveToken SaveToken
	// Decode from base64
	var b, err = base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return nil, err
	}
	left, err := saveToken.UnmarshalMsg(b)
	if err != nil {
		return nil, err
	}
	if len(left) > 0 {
		return nil, errors.New("left over bytes")
	}
	// Get the data
	tokenToSet.AccessToken = saveToken.AccessToken
	tokenToSet.RefreshToken = saveToken.RefreshToken
	tokenToSet.LastUpdate = saveToken.LastUpdate

	return tokenToSet, nil
}

// Delete the token cookie
func DeleteTokenCookie() {
	jsext.DeleteCookie(JSEXT_token)
}
