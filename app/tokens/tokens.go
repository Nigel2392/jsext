package tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Nigel2392/jsext/requester"
)

var AuthToken *Token

type TokenURLs struct {
	LoginURL    string
	RegisterURL string
	LogoutURL   string
	RefreshURL  string
}

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
	stopChan             chan bool
}

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
		Data:     make(map[string]interface{}),
		stopChan: make(chan bool),
	}
	return t
}

func (t *Token) SetURLs(urls TokenURLs) {
	t.URLs = urls
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.LastUpdate.Add(t.AccessTimeout))
}

func (t *Token) IsRefreshExpired() bool {
	return time.Now().After(t.LastUpdate.Add(t.RefreshTimeout))
}

func (t *Token) ExpiredIn() time.Duration {
	timeout := t.LastUpdate.Add(t.AccessTimeout)
	return time.Until(timeout)
}

func (t *Token) RefreshExpiredIn() time.Duration {
	timeout := t.LastUpdate.Add(t.RefreshTimeout)
	return time.Until(timeout)
}

func (t *Token) Update() error {
	var client = requester.NewAPIClient()
	var data = map[string]string{
		"refresh": t.RefreshToken,
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
	return <-errChan
}

func (t *Token) Client() *requester.APIClient {
	var client = requester.NewAPIClient()
	client.OnError(func(err error) bool {
		println(err.Error())
		return true
	})
	client = client.WithHeaders(map[string][]string{
		"Authorization": {"Bearer " + t.AccessToken},
	})
	return client
}

func (t *Token) sendDataGetToken(data map[string]string, url string) error {
	var client = requester.NewAPIClient()
	client = client.Post(url)
	client.OnError(func(err error) bool {
		println(err.Error())
		return false
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
		if datamap[t.errorMessageName] != nil {
			errChan <- errors.New(datamap["detail"].(string))
			return
		}
		t.AccessToken = datamap[t.AccessTokenVariable].(string)
		t.RefreshToken = datamap[t.RefreshTokenVariable].(string)
		delete(datamap, t.AccessTokenVariable)
		delete(datamap, t.RefreshTokenVariable)
		println(fmt.Sprintf("%v", datamap))
		t.Data = datamap
		t.LastUpdate = time.Now()
		t.updateManager()
		errChan <- nil
	})
	return <-errChan
}

func (t *Token) Login(loginData map[string]string) error {
	return t.sendDataGetToken(loginData, t.URLs.LoginURL)
}

func (t *Token) Register(registerData map[string]string) error {
	return t.sendDataGetToken(registerData, t.URLs.RegisterURL)
}

func (t *Token) Logout() error {
	var refresh = t.RefreshToken
	var client = t.Client()
	client = client.Post(t.URLs.LogoutURL)
	client.WithData(map[string]string{
		t.RefreshTokenVariable: refresh,
	}, requester.JSON)
	var errChan = make(chan error)
	var respMap map[string]any
	client.DoDecodeTo(&respMap, requester.JSON, func(r *http.Response, s any) {
		var err, ok = respMap[t.errorMessageName].(string)
		if ok {
			errChan <- errors.New(err)
			return
		}
		errChan <- nil
	})
	var err = <-errChan
	t.Reset()
	return err
}

func (t *Token) updateManager() {
	go func() {
		for {
			switch {
			case <-t.stopChan:
				return
			default:
				var tim = t.ExpiredIn()
				if tim <= 3*time.Minute {
					t.Update()
				}
				tim = t.ExpiredIn()
				time.Sleep(tim - 3*time.Minute)
			}
		}
	}()
}

func (t *Token) stopManager() {
	t.stopChan <- true
	close(t.stopChan)
}

func (t *Token) Reset() *Token {
	var urls = t.URLs
	t.stopManager()
	t = NewToken(t.RefreshTimeout, t.AccessTimeout, t.AccessTokenVariable, t.RefreshTokenVariable, t.errorMessageName)
	t.SetURLs(urls)
	return t
}
