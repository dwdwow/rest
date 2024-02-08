package rest

import "net/http"

type Header map[string]string

type Response struct {
	Raw     *http.Response
	Request *http.Request
	Data    []byte
}
