//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package requester_test

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/Nigel2392/jsext/requester"
)

const ADDRESS = "127.0.0.1:8080"

type Data struct {
	Status  string `json:"status" xml:"status"`
	Message string `json:"message" xml:"message"`
}

func init() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]interface{}{
			"status":  "ok",
			"message": "Hello, World!",
		}
		w.Header().Set("Content-Type", "application/json")
		var err = json.NewEncoder(w).Encode(data)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello, World!</h1>"))
	})

	http.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		var data = Data{
			Status:  "ok",
			Message: "Hello, World!",
		}
		var err = xml.NewEncoder(w).Encode(data)
		if err != nil {
			panic(err)
		}
	})

	go func() {
		http.ListenAndServe(ADDRESS, nil)
	}()

}

func TestGetJSON(t *testing.T) {
	ch := make(chan struct{})
	var c = requester.NewAPIClient()
	var data = make(map[string]interface{})
	c.Get("http://"+ADDRESS+"/json").DoDecodeTo(&data, requester.JSON, func(resp *http.Response, strct interface{}) {
		if resp.StatusCode != http.StatusOK {
			t.Error("Status code is not 200")
		}
		if data["status"] != "ok" {
			t.Error("Status is not ok")
		}
		if data["message"] != "Hello, World!" {
			t.Error("Message is not 'Hello, World!'")
		}
		ch <- struct{}{}
	})
	<-ch
}

func TestGetXML(t *testing.T) {
	ch := make(chan struct{})
	var c = requester.NewAPIClient()
	var data = Data{}
	c.Get("http://"+ADDRESS+"/xml").DoDecodeTo(&data, requester.XML, func(resp *http.Response, strct interface{}) {
		if resp.StatusCode != http.StatusOK {
			t.Error("Status code is not 200")
		}
		if data.Status != "ok" {
			t.Error("Status is not ok")
		}
		if data.Message != "Hello, World!" {
			t.Error("Message is not 'Hello, World!'")
		}
		ch <- struct{}{}
	})
	<-ch
}

func TestGetHTML(t *testing.T) {
	var c = requester.NewAPIClient()
	ch := make(chan struct{})
	c.Get("http://" + ADDRESS + "/html").Do(func(resp *http.Response) {
		if resp.StatusCode != http.StatusOK {
			t.Error("Status code is not 200")
		}
		if resp.Header.Get("Content-Type") != "text/html" {
			t.Error("Content-Type is not text/html")
		}
		data := make([]byte, resp.ContentLength)
		resp.Body.Read(data)
		if string(data) != "<h1>Hello, World!</h1>" {
			t.Error("Body is not '<h1>Hello, World!</h1>'")
		}
		ch <- struct{}{}
	})
	<-ch
}
