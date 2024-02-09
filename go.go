package rest

import (
	"bytes"
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

const minStatusCode = 399

type Header map[string]string

// Body must be a map, struct and string.
type Body any

type Response struct {
	Raw     *http.Response
	Request *http.Request
	Data    []byte

	RespBodyCloseErr error
}

func Go(method Method, url string, header Header, body io.Reader, result any) (*Response, error) {
	return GoWithClient(http.DefaultClient, method, url, header, body, result)
}

func GoWithClient(client *http.Client, method Method, url string, header Header, body Body, result any) (*Response, error) {
	if client == nil {
		return nil, fmt.Errorf("rest: client is nil")
	}

	bd, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("reset: marshal body to bytes, err: %w", err)
	}

	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(bd))
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

	r.Raw, err = client.Do(req)
	if err != nil {
		return r, err
	}
	defer func() {
		r.RespBodyCloseErr = r.Raw.Body.Close()
	}()

	if r.Raw.StatusCode > minStatusCode {
		return r, fmt.Errorf("rest: status %v > %v", r.Raw.StatusCode, minStatusCode)
	}

	r.Data, err = io.ReadAll(r.Raw.Body)
	if err != nil {
		return r, fmt.Errorf("rest: can not read response body, err: %w", err)
	}

	if result != nil && r.Data != nil {
		err = json.Unmarshal(r.Data, result)
		if err != nil {
			return r, fmt.Errorf("rest: can not unmarshal response data, err: %w", err)
		}
	}

	return r, nil
}

func Get(url string, header Header, result any) (*Response, error) {
	return Go(GET, url, header, nil, result)
}
