package requests

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendRequestGetResp(method string, url string, reqBody *bytes.Buffer, headersMap map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(
		method,
		url,
		reqBody,
	)
	if err != nil {
		return nil, err
	}

	for key, value := range headersMap {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed) {
		return nil, fmt.Errorf("response status: %d", resp.StatusCode)
	}

	return resp, nil
}
