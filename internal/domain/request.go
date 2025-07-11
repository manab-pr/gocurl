package domain

import "net/http"

type Request struct {
	Method  string
	URL     string
	Body    string
	Headers http.Header
}
