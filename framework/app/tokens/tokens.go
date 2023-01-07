//go:build js && wasm
// +build js,wasm

package tokens

import (
	"reflect"
	"time"
)

var AuthToken *Token

// URLs for the token to request to.
type TokenURLs struct {
	LoginURL    string
	RegisterURL string
	LogoutURL   string
	RefreshURL  string
}

// Authentication token struct
type Token struct {
	AccessToken          string
	RefreshToken         string
	LastUpdate           time.Time
	RefreshTimeout       time.Duration
	AccessTimeout        time.Duration
	URLs                 TokenURLs
	AccessTokenVariable  string
	RefreshTokenVariable string
	errorMessageName     string
	Data                 map[string]interface{}
	onUpdate             func(t *Token)
	onReset              func()
	onInit               func(t *Token)
	onUpdateErr          func(err error)
	needsUpdateChan      chan bool
}

// Get a new token
func NewToken(RefreshTimeout, AccessTimeout time.Duration, AccessVar, RefreshVar, errorMessageName string) *Token {
	var t = &Token{
		RefreshTimeout:       RefreshTimeout,
		AccessTimeout:        AccessTimeout,
		AccessTokenVariable:  AccessVar,
		RefreshTokenVariable: RefreshVar,
		errorMessageName:     errorMessageName,
		URLs: TokenURLs{
			RegisterURL: "http://127.0.0.1:8000/api/auth/register/",
			LogoutURL:   "http://127.0.0.1:8000/api/auth/logout/",
			LoginURL:    "http://127.0.0.1:8000/api/auth/login/",
			RefreshURL:  "http://127.0.0.1:8000/api/auth/refresh/",
		},
		Data: make(map[string]interface{}),
	}
	return t
}

// Set the token URLs
func (t *Token) SetURLs(urls TokenURLs) {
	t.URLs = urls
}

// Callback for when the token gets updated.
func (t *Token) OnUpdate(f func(t *Token)) {
	t.onUpdate = f
}

// Callback for when the token gets reset.
func (t *Token) OnReset(f func()) {
	t.onReset = f
}

// Callback for when the token gets initialized.
func (t *Token) OnInit(f func(t *Token)) {
	t.onInit = f
}

// Callback when an error occurs updating the token.
func (t *Token) OnUpdateError(f func(err error)) {
	t.onUpdateErr = f
}

// Needs update returns a channel which will send a bool when the token needs to be updated.
func (t *Token) NeedsUpdate() <-chan bool {
	if t.needsUpdateChan == nil {
		t.needsUpdateChan = make(chan bool)
	}
	go func() {
		for {
			// Check if the token needs to be updated.
			if t.ShouldUpdate() {
				t.needsUpdateChan <- true
				return
			}
			time.Sleep(t.ExpiredIn() - time.Duration(t.ExpiredIn()/10))
		}
	}()
	return t.needsUpdateChan
}

func (t *Token) ShouldUpdate() bool {
	if t.AccessToken == "" ||
		t.RefreshToken == "" ||
		t.RefreshTimeout == 0 ||
		t.AccessTimeout == 0 ||

		t.IsRefreshExpired() {
		return false
	} else if t.IsExpired() {
		return true
	} else if t.LastUpdate.Add(t.AccessTimeout - time.Duration(t.AccessTimeout/10)).Before(time.Now()) {
		return true
	}
	return false
}

// Check if access token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.LastUpdate.Add(t.AccessTimeout))
}

// Check if refresh token is expired
func (t *Token) IsRefreshExpired() bool {
	return time.Now().After(t.LastUpdate.Add(t.RefreshTimeout))
}

// Get when the access token will expire
func (t *Token) ExpiredIn() time.Duration {
	timeout := t.LastUpdate.Add(t.AccessTimeout)
	return time.Until(timeout)
}

// Get when the refresh token will expire
func (t *Token) RefreshExpiredIn() time.Duration {
	timeout := t.LastUpdate.Add(t.RefreshTimeout)
	return time.Until(timeout)
}

// Short hand for getting the token data
func (t *Token) setToken(access, refresh string, lastUpdate time.Time) {
	t.AccessToken = access
	t.RefreshToken = refresh
	t.LastUpdate = lastUpdate
	t.RunManager()
}

// Run the token update manager.
// This will automatically update the token every AccessTimeout - 10%
// Automatically stops the manager if an error occurs when updating the token.
func (t *Token) RunManager() {
	t.StopManager()
	go func() {
		for <-t.NeedsUpdate() {
			var err = t.Update()
			if err != nil {
				if t.onUpdateErr != nil {
					t.onUpdateErr(err)
					t.StopManager()
					return
				} else {
					panic(err)
				}
			}
		}
	}()
}

// Stop the token update manager.
func (t *Token) StopManager() {
	if t.needsUpdateChan != nil {
		close(t.needsUpdateChan)
		t.needsUpdateChan = nil
	}
}

// Reset the token.
// Essentially creates a new token, transfers all nescessaary data and overwrites the old one.
func (t *Token) Reset() *Token {
	var urls = t.URLs
	if t.onReset != nil {
		t.onReset()
	}
	t.StopManager()
	DeleteTokenCookie()
	var newt = NewToken(t.RefreshTimeout, t.AccessTimeout, t.AccessTokenVariable, t.RefreshTokenVariable, t.errorMessageName)
	newt.OnInit(t.onInit)
	newt.OnUpdate(t.onUpdate)
	newt.OnReset(t.onReset)
	newt.SetURLs(urls)
	reflect.ValueOf(t).Elem().Set(reflect.ValueOf(newt).Elem())
	return t
}
