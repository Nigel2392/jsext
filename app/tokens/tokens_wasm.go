package tokens

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"time"

	"github.com/Nigel2392/jsext"
)

func SetTokenCookie(token *Token) error {
	var AccessToken = token.AccessToken
	var RefreshToken = token.RefreshToken
	var LastUpdate = token.LastUpdate

	datamap := make(map[string]interface{})
	datamap["AccessToken"] = AccessToken
	datamap["RefreshToken"] = RefreshToken
	datamap["LastUpdate"] = LastUpdate

	var b bytes.Buffer
	// Gob the token
	var enc = gob.NewEncoder(&b)
	enc.Encode(datamap)
	// Encode to base64
	var cookie = base64.RawURLEncoding.EncodeToString(b.Bytes())
	// Set the cookie
	return jsext.SetCookie("token", cookie, 1*3600*24)
}

func GetTokenCookie(tokenToSet *Token) (*Token, error) {
	var cookie = jsext.GetCookie("token")
	if cookie == "" {
		return nil, errors.New("No token cookie")
	}
	// Decode from base64
	var b, err = base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return nil, err
	}
	// Gob the token
	var dec = gob.NewDecoder(bytes.NewBuffer(b))
	var datamap map[string]interface{}
	dec.Decode(&datamap)
	// Get the data
	var AccessToken = datamap["AccessToken"].(string)
	var RefreshToken = datamap["RefreshToken"].(string)
	var LastUpdate = datamap["LastUpdate"].(time.Time)
	// Create the token
	tokenToSet.AccessToken = AccessToken
	tokenToSet.RefreshToken = RefreshToken
	tokenToSet.LastUpdate = LastUpdate
	return tokenToSet, nil
}
