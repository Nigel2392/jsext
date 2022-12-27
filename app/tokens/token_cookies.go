package tokens

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/Nigel2392/jsext"
)

const JSEXT_token = "JSEXT_token"

// Set the token cookie from a token.
func SetTokenCookie(token *Token) error {
	var AccessToken = token.AccessToken
	var RefreshToken = token.RefreshToken
	var LastUpdate = token.LastUpdate

	datamap := make(map[string]interface{})
	datamap["AccessToken"] = AccessToken
	datamap["RefreshToken"] = RefreshToken
	datamap["LastUpdate"] = LastUpdate
	// Json the token
	var b bytes.Buffer
	var err = json.NewEncoder(&b).Encode(datamap)
	if err != nil {
		return err
	}
	// Encode to base64
	var cookie = base64.RawURLEncoding.EncodeToString(b.Bytes())
	// Set the cookie
	return jsext.SetCookie(JSEXT_token, cookie, time.Second*3600*24)
}

// Get the token cookie
func GetTokenCookie(tokenToSet *Token) (*Token, error) {
	var cookie = jsext.GetCookie(JSEXT_token)
	if cookie == "" {
		//lint:ignore ST1005 Error strings should not be capitalized
		return nil, errors.New("No token cookie")
	}
	// Decode from base64
	var b, err = base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return nil, err
	}
	// Json the token
	var datamap map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&datamap)
	if err != nil {
		return nil, err
	}
	// Get the data
	AccessToken, ok := datamap["AccessToken"].(string)
	if !ok {
		//lint:ignore ST1005 Error strings should not be capitalized
		return nil, errors.New("No cookie access token")
	}
	RefreshToken, ok := datamap["RefreshToken"].(string)
	if !ok {
		//lint:ignore ST1005 Error strings should not be capitalized
		return nil, errors.New("No cookie refresh token")
	}
	LastUpdate, ok := datamap["LastUpdate"].(string)
	if !ok {
		return nil, errors.New("Token cookie time could not be parsed")
	}
	LastUpdateParsed, err := time.Parse(time.RFC3339, LastUpdate)
	if err != nil {
		return nil, err
	}
	// Create the token
	tokenToSet.AccessToken = AccessToken
	tokenToSet.RefreshToken = RefreshToken
	tokenToSet.LastUpdate = LastUpdateParsed
	if tokenToSet.IsExpired() && !tokenToSet.IsRefreshExpired() {
		tokenToSet.Update()
	} else if tokenToSet.IsRefreshExpired() {
		return nil, errors.New("Token cookie expired")
	} else if tokenToSet.ExpiredIn() < tokenToSet.AccessTimeout-time.Duration(tokenToSet.AccessTimeout.Seconds()/10) {
		tokenToSet.Update()
	}
	return tokenToSet, nil
}

// Delete the token cookie
func DeleteTokenCookie() {
	jsext.DeleteCookie(JSEXT_token)
}
