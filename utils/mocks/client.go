package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MockClient is the mock client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

// SuccessFromFile creates simple `Do` func which responds with 200 and the contents of given file
func SuccessFromFile(filename string) func(req *http.Request) (*http.Response, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	r := ioutil.NopCloser(bytes.NewReader(contents))
	return func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
}
