//go:build tinygo && js && wasm
// +build tinygo,js,wasm

package tokens

import (
	"errors"
	"time"

	"github.com/Nigel2392/jsext/framework/requester"
	"github.com/tidwall/gjson"
)

// Make an api call to the refresh URL, update both the access and refresh tokens.
func (t *Token) Update() error {
	var client = requester.NewAPIClient()
	var data = map[string]any{
		t.RefreshTokenVariable: t.RefreshToken,
	}
	client = client.Post(t.URLs.RefreshURL).WithData(data, requester.JSON)
	var resp, err = client.Do()
	if err != nil {
		return err
	}
	var datamap, ok = gjson.ParseBytes(resp.Body).Value().(map[string]interface{})
	if !ok {
		return errors.New("could not parse response")
	}
	if err, ok := datamap[t.errorMessageName]; ok {
		return errors.New(err.(string))
	}
	AccessToken, ok := datamap[t.AccessTokenVariable].(string)
	if !ok {
		return errors.New("could not parse response")
	}
	RefreshToken, ok := datamap[t.RefreshTokenVariable].(string)
	if !ok {
		return errors.New("could not parse response")
	}
	var LastUpdate = time.Now()
	t.AccessToken = AccessToken
	t.RefreshToken = RefreshToken
	t.LastUpdate = LastUpdate
	if t.onUpdate != nil {
		t.onUpdate(t)
	}
	return nil
}

// Token's client, when authenticated it will automatically set the Authorization header.
func (t *Token) Client() *requester.APIClient {
	var client = requester.NewAPIClient()
	if t.AccessToken == "" {
		return client
	}
	client = client.WithHeaders(map[string]string{
		"Authorization": "Bearer " + t.AccessToken,
	})
	return client
}

// Send data to an API endpoint, get both access and refresh tokens.
func (t *Token) sendDataGetToken(data map[string]any, url string) error {
	var client = requester.NewAPIClient()
	client = client.Post(url)
	client.WithData(data, requester.JSON)
	var resp, err = client.Do()
	if err != nil {
		return err
	}
	var datamap, ok = gjson.ParseBytes(resp.Body).Value().(map[string]interface{})
	if !ok {
		return errors.New("could not parse response")
	}
	if err, ok := datamap[t.errorMessageName]; ok {
		return errors.New(err.(string))
	}
	var AccessToken = datamap[t.AccessTokenVariable].(string)
	var RefreshToken = datamap[t.RefreshTokenVariable].(string)
	var LastUpdate = time.Now()
	delete(datamap, t.AccessTokenVariable)
	delete(datamap, t.RefreshTokenVariable)
	t.Data = datamap
	t.setToken(AccessToken, RefreshToken, LastUpdate)
	if t.onInit != nil {
		t.onInit(t)
	}
	return nil
}

// Login with the appropriate data, get both access and refresh tokens.
func (t *Token) Login(loginData map[string]string) error {
	newMap := make(map[string]any, len(loginData))
	for k, v := range loginData {
		newMap[k] = v
	}
	return t.sendDataGetToken(newMap, t.URLs.LoginURL)
}

// Register with the appropriate data, get both access and refresh tokens.
func (t *Token) Register(registerData map[string]string) error {
	newMap := make(map[string]any, len(registerData))
	for k, v := range registerData {
		newMap[k] = v
	}
	return t.sendDataGetToken(newMap, t.URLs.RegisterURL)
}

// Logout with the refresh token.
func (t *Token) Logout() error {
	if t.AccessToken == "" || t.RefreshToken == "" || t.URLs.LogoutURL == "" {
		//lint:ignore ST1005 Error strings should not be capitalized
		return errors.New("Already logged out")
	}
	var client = t.Client()
	client = client.Post(t.URLs.LogoutURL)
	client.WithData(map[string]any{
		t.RefreshTokenVariable: t.RefreshToken,
	}, requester.JSON)
	var resp, err = client.Do()
	if err != nil {
		return err
	}
	var respMap, ok = gjson.ParseBytes(resp.Body).Value().(map[string]interface{})
	if !ok {
		return errors.New("could not parse response")
	}
	if err, ok := respMap[t.errorMessageName]; ok {
		return errors.New(err.(string))
	}
	t.Reset()
	return nil
}
