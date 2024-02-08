package rest

import "net/http"

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
}
