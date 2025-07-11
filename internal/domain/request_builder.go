package domain

import (
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

type RequestBuilder struct {
	request *Request
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{&Request{Headers: make(http.Header)}}}

func (b *RequestBuilder) WithMethod(method string) *RequestBuilder {
	b.request.Method = strings.ToUpper(method)
	return b
}

func (b *RequestBuilder) WithURL(url string) *RequestBuilder {
	b.request.URL = url
	return b
}

func (b *RequestBuilder) WithBody(body string) *RequestBuilder {
	if strings.HasPrefix(body, "@") {
		filename := strings.TrimPrefix(body, "@")
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			color.Red("❌ Failed to read file %s: %v\n", filename, err)
			return b
		}
		b.request.Body = string(fileContent)
	} else {
		b.request.Body = body
	}
	return b
}

func (b *RequestBuilder) WithHeaders(headers []string) *RequestBuilder {
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			b.request.Headers.Add(key, val)
		}
	}
	return b
}

func (b *RequestBuilder) Build() *Request {
	if b.request.Body != "" && b.request.Headers.Get("Content-Type") == "" {
		b.request.Headers.Set("Content-Type", "application/json")
	}
	return b.request
}

