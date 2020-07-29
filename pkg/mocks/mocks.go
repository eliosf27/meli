package mocks

import (
	"github.com/jarcoal/httpmock"
	"net/http"
)

type (
	HttpMock struct{}
)

func NewHttpMock() HttpMock {
	return HttpMock{}
}

// Activate starts the mock environment
func (HttpMock) Activate() {
	httpmock.Activate()
}

// DeactivateAndReset shuts down the mock environment
func (HttpMock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

// httpMockResponse adds a new responder, associated with a given HTTP method and URL with a mock data interface
func (HttpMock) httpMockResponse(method string, url string, data interface{}) {
	httpmock.RegisterResponder(method, url,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, data)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

// Post adds a new responder, associated with the HTTP POST method and URL with a mock data interface
func (mockHttp HttpMock) Post(url string, data interface{}) {
	mockHttp.httpMockResponse("POST", url, data)
}

// Post adds a new responder, associated with the HTTP GET method and URL with a mock data interface
func (mockHttp HttpMock) Get(url string, data interface{}) {
	mockHttp.httpMockResponse("GET", url, data)
}
