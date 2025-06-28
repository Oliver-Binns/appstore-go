package mocknetworking

import (
	"bytes"
	"io"
	"net/http"
)

type MockHTTPClient struct {
	Requests  []*http.Request
	Responses []MockHTTPResponse
}

type MockHTTPResponse struct {
	StatusCode *int
	Body       string
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.Requests = append(c.Requests, req)
	response := c.Responses[0]
	c.Responses = c.Responses[1:]

	responseBody := io.NopCloser(bytes.NewReader([]byte(response.Body)))

	status := http.StatusOK
	if response.StatusCode != nil {
		status = *response.StatusCode
	}

	return &http.Response{
		StatusCode: status,
		Body:       responseBody,
	}, nil
}
