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

func (HttpMock) Activate() {
	httpmock.Activate()
}

func (HttpMock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

func (HttpMock) HttpMockResponse(method string, url string, data interface{}) {
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

func (mockHttp HttpMock) Post(url string, data interface{}) {
	mockHttp.HttpMockResponse("POST", url, data)
}

func (mockHttp HttpMock) Get(url string, data interface{}) {
	mockHttp.HttpMockResponse("GET", url, data)
}

func (mockHttp HttpMock) Calls(path string) int {
	calls := httpmock.GetCallCountInfo()
	if res, ok := calls[path]; ok {
		return res
	}

	return 0
}
