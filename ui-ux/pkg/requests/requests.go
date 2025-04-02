package requests

import (
	"bytes"
	"net/http"
)

func SendRequest(method, url string, reqBody *bytes.Buffer, headers map[string]string) (*http.Response, error) {
	if reqBody == nil {
		reqBody = &bytes.Buffer{} // Empty buffer to prevent nil dereference
	}

	req, err := http.NewRequest(
		method,
		url,
		reqBody,
	)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
