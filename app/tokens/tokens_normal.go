//go:build !tinygo && js && wasm
// +build !tinygo,js,wasm

package tokens

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Nigel2392/jsext/requester"
)

// Make an api call to the refresh URL, update both the access and refresh tokens.
func (t *Token) Update() error {
	var client = requester.NewAPIClient()
	var data = map[string]any{
		t.RefreshTokenVariable: t.RefreshToken,
	}
	client = client.Post(t.URLs.RefreshURL).WithData(data, requester.JSON)
	var datamap map[string]string
	var errChan = make(chan error)
	client.Do(func(resp *http.Response) {
		var err = json.NewDecoder(resp.Body).Decode(&datamap)
		if err != nil {
			errChan <- err
			return
		}
		t.AccessToken = datamap[t.AccessTokenVariable]
		t.RefreshToken = datamap[t.RefreshTokenVariable]
		t.LastUpdate = time.Now()
		errChan <- nil
	})
	var err = <-errChan
	if err != nil {
		return err
	}
	if t.onUpdate != nil {
		t.onUpdate(t)
	}
	return nil
}

// Token's client, when authenticated it will automatically set the Authorization header.
func (t *Token) Client() *requester.APIClient {
	var client = requester.NewAPIClient()
	client.OnError(func(err error) bool {
		println(err.Error())
		return true
	})
	if t.AccessToken == "" {
		return client
	}
	client = client.WithHeaders(map[string][]string{
		"Authorization": {"Bearer " + t.AccessToken},
	})
	return client
}

// Send data to an API endpoint, get both access and refresh tokens.
func (t *Token) sendDataGetToken(data map[string]any, url string) error {
	var client = requester.NewAPIClient()
	client = client.Post(url)
	client.OnError(func(err error) bool {
		println(err.Error())
		return true
	})
	client.WithData(data, requester.JSON)
	var datamap map[string]interface{}
	var errChan = make(chan error)
	client.Do(func(resp *http.Response) {
		var err = json.NewDecoder(resp.Body).Decode(&datamap)
		if err != nil {
			errChan <- err
			return
		}
		if err, ok := datamap[t.errorMessageName]; ok {
			errChan <- errors.New(err.(string))
			return
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
		errChan <- nil
	})
	return <-errChan
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
	var errChan = make(chan error)
	var respMap map[string]any
	client.DoDecodeTo(&respMap, requester.JSON, func(r *http.Response, s any) {
		var err, ok = respMap[t.errorMessageName]
		if ok {
			errChan <- errors.New(err.(string))
			return
		}
		errChan <- nil
	})
	var err = <-errChan
	t.Reset()
	return err
}

// Decode the token into its parts
func (t *Token) JWTDecode() (JWTToken, JWTToken, error) {
	access, err := DecodeToken(t.AccessToken)
	if err != nil {
		return JWTToken{}, JWTToken{}, err
	}
	refresh, err := DecodeToken(t.RefreshToken)
	if err != nil {
		return JWTToken{}, JWTToken{}, err
	}
	return access, refresh, nil
}
