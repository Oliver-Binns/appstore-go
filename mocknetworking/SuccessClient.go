package mocknetworking

func MockHTTPClientWith200Response(body string) *MockHTTPClient {
	return &MockHTTPClient{
		Responses: []MockHTTPResponse{
			{
				Body: body,
			},
		},
	}
}

func MockHTTPClientWithSingleResponse(statusCode int, body string) *MockHTTPClient {
	return &MockHTTPClient{
		Responses: []MockHTTPResponse{
			{
				StatusCode: &statusCode,
				Body:       body,
			},
		},
	}
}
