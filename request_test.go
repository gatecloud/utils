package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Coffee struct {
	Size string
	Type string
}

func TestRequestJson(t *testing.T) {
	wantStatusCode := 200
	want := Coffee{
		Size: "Medium",
		Type: "Flat White",
	}

	// Create mock http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		c := Coffee{
			Size: "Medium",
			Type: "Flat White",
		}

		bytes, err := json.Marshal(c)
		if err != nil {
			t.Error(err)
		}
		io.WriteString(w, string(bytes))
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	var got Coffee
	request := Request{
		Method: "GET",
		URL:    ts.URL,
		Header: make(http.Header),
		Object: &got,
	}

	resp, err := request.Do()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, wantStatusCode, resp.StatusCode)
	assert.Equal(t, want, got)
}

func TestRequestBinary(t *testing.T) {
	wantStatusCode := 200
	want, err := ioutil.ReadFile("golang.png")
	if err != nil {
		t.Error(err)
	}

	// Create mock http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile("golang.png")
		if err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "image/png")
		io.WriteString(w, string(bytes))
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	var got []byte
	request := Request{
		Method: "GET",
		URL:    ts.URL,
		Header: make(http.Header),
		Object: &got,
	}

	resp, err := request.Do()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, wantStatusCode, resp.StatusCode)
	assert.Equal(t, want, got)
}
