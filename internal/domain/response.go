package domain

import "net/http"

type Response struct {
	Status  string
	Headers http.Header
	Body    []byte
}
