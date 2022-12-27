package tokens

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/Nigel2392/jsext/requester"
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
	ticker               *time.Ticker
	onUpdate             func(t *Token)
	onReset              func()
	onInit               func(t *Token)
	onUpdateErr          func(err error)
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
	c := make(chan bool)
	go func() {
		for {
			if t.IsExpired() || t.ExpiredIn() < t.AccessTimeout/10 {
				c <- true
				return
			}
			time.Sleep(t.AccessTimeout / 20)
		}
	}()
	return c
}

func (t *Token) ShouldUpdate() bool {
	if t.IsExpired() && !t.IsRefreshExpired() {
		return true
	} else if t.IsRefreshExpired() {
		return false
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

// Make an api call to the refresh URL, update both the access and refresh tokens.
func (t *Token) Update() error {
	var client = requester.NewAPIClient()
	var data = map[string]string{
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
func (t *Token) sendDataGetToken(data map[string]string, url string) error {
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
	return t.sendDataGetToken(loginData, t.URLs.LoginURL)
}

// Register with the appropriate data, get both access and refresh tokens.
func (t *Token) Register(registerData map[string]string) error {
	return t.sendDataGetToken(registerData, t.URLs.RegisterURL)
}

// Logout with the refresh token.
func (t *Token) Logout() error {
	if t.AccessToken == "" || t.RefreshToken == "" || t.URLs.LogoutURL == "" {
		//lint:ignore ST1005 Error strings should not be capitalized
		return errors.New("Already logged out")
	}
	var client = t.Client()
	client = client.Post(t.URLs.LogoutURL)
	client.WithData(map[string]string{
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

// Run the token update manager.
// This will automatically update the token every AccessTimeout - 10%
func (t *Token) RunManager() {
	t.StopManager()
	t.ticker = time.NewTicker(t.AccessTimeout - time.Duration(t.AccessTimeout/10))
	// Takes a bool to update the token immediately.
	// This might be nescessary if the token expires before the first update.
	// Note: This will panic if the token update fails and no onUpdateErr handler is set.

	//if update {
	//	var err = t.Update()
	//	if err != nil {
	//		if t.onUpdateErr != nil {
	//			t.onUpdateErr(err)
	//		} else {
	//			panic(err)
	//		}
	//	}
	//}
	go func() {
		for {
			select {
			case <-t.ticker.C:
				if t.AccessToken == "" || t.RefreshToken == "" {
					continue
				}
				var err = t.Update()
				if err != nil {
					if t.onUpdateErr != nil {
						t.onUpdateErr(err)
					} else {
						t.ticker.Stop()
						panic(err)
					}
				}
			case <-t.NeedsUpdate():
				var err = t.Update()
				if err != nil {
					if t.onUpdateErr != nil {
						t.onUpdateErr(err)
					} else {
						t.ticker.Stop()
						panic(err)
					}
				}
			}
		}
	}()
}

// Stop the token update manager.
func (t *Token) StopManager() {
	if t.ticker != nil {
		t.ticker.Stop()
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
