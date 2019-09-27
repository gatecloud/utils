package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Request stores the data that http request needs
// and get the corresponding object data and header
//
// Header must be created before using
// e.g. RequestHeader = make(http.Header)
// Body allows to pass an entity of JSON model
// Object should be pointer type
type Request struct {
	Method string
	URL    string
	Header http.Header
	Body   interface{}
	Object interface{}
	Retry  int
}

type Response struct {
	Status     string
	StatusCode int
	Header     http.Header
}

// Request
func (r *Request) Do() (Response, error) {
	var (
		err      error
		response Response
		reader   io.Reader
	)
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if r.Header == nil {
		return Response{StatusCode: 400}, errors.New("Request Header is not initialized")
	}

	method := strings.ToUpper(r.Method)
	if method != "GET" && method != "HEAD" && method != "OPTION" {
		if !reflect.ValueOf(r.Body).IsValid() {
			return Response{StatusCode: 400}, errors.New("Request Body can't be empty")
		}

		data, err := json.Marshal(r.Body)
		if err != nil {
			return Response{Status: "Request Marshal Error", StatusCode: 500}, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(r.Method, r.URL, reader)
	if err != nil {
		return Response{Status: "Request 500 Internal Server Error", StatusCode: 500}, err
	}
	req = req.WithContext(c)
	req.Header = r.Header
	if req.Header.Get("Content-Type") == "" || req.Header.Get("content-type") == "" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Connection", "close")
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		response.Status = resp.Status
		response.StatusCode = resp.StatusCode
		response.Header = make(http.Header)
		response.Header = resp.Header
	}
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		if r.Retry < 2 {
			return response, errors.New(resp.Status)
		}
		r.Retry--
		return r.Do()
	}

	if resp.StatusCode == 204 {
		return response, nil
	}

	resBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return response, err
	}

	if (resp.StatusCode >= 400 && resp.StatusCode < 500) || resp.StatusCode < 200 {
		stringBody := string(resBody)
		return response, errors.New(stringBody)
	}

	if !reflect.ValueOf(r.Object).IsValid() {
		return Response{StatusCode: 400}, errors.New("Request Object can't be nil")
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "image") {
		obj, ok := r.Object.(*[]byte)
		if !ok {
			return Response{Status: "Request Assert Error", StatusCode: 500}, err
		}
		*obj = resBody
		r.Object = obj

	} else {
		err = json.Unmarshal(resBody, r.Object)
		if err != nil {
			return Response{Status: "Request Unmarshal Error", StatusCode: 500}, err
		}
	}
	return response, nil
}
