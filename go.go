package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type Header map[string]string

type Response struct {
	Raw     *http.Response
	Request *http.Request
	Data    []byte

	RespBodyCloseErr error
}

func GoWithClient(client *http.Client, method Method, url string, header Header, body io.Reader, result any) (*Response, error) {
	if client == nil {
		return nil, fmt.Errorf("rest: client is nil")
	}

	req, err := http.NewRequest(string(method), url, body)
	if err != nil {
		return nil, fmt.Errorf("rest: can not new request, err: %w", err)
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	r := &Response{
		Request: req,
	}

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer func() {
		r.RespBodyCloseErr = resp.Body.Close()
	}()

	r.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		return r, fmt.Errorf("rest: can not read response body, err: %w", err)
	}

	r.Raw = resp

	if result != nil && r.Data != nil {
		err = json.Unmarshal(r.Data, result)
		if err != nil {
			return r, fmt.Errorf("rest: can not unmarshal response data, err: %w", err)
		}
	}

	return r, nil
}
