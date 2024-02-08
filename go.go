package rest

import "net/http"

type Response struct {
	Raw     *http.Response
	Request *http.Request
	Data    []byte
}
