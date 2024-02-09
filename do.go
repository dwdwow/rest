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

// Body must be a map, struct or string.
type Body any

// Result must be a map, or a pointer to a struct or string
type Result any

type Response struct {
	Raw     *http.Response
	Request *http.Request
	Data    []byte

	RespBodyCloseErr error
}

func Go(method Method, url string, header Header, body Body, result Result) (*Response, error) {
	return GoWithClient(http.DefaultClient, method, url, header, body, result)
}

func GoWithClient(client *http.Client, method Method, url string, header Header, body Body, result Result) (*Response, error) {
	if client == nil {
		return nil, fmt.Errorf("rest: client is nil")
	}

	bd, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("reset: marshal body to bytes, err: %w", err)
	}

	req, err := http.NewRequest(string(method), url, bytes.NewReader(bd))
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

func Get(url string, header Header, result Result) (*Response, error) {
	return Go(GET, url, header, nil, result)
}

func Post(url string, header Header, body Body, result Result) (*Response, error) {
	return Go(POST, url, header, body, result)
}

func Put(url string, header Header, body Body, result Result) (*Response, error) {
	return Go(PUT, url, header, body, result)
}

func Delete(url string, header Header, result Result) (*Response, error) {
	return Go(DELETE, url, header, nil, result)
}
