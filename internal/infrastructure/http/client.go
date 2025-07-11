package http

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manab-pr/gocurl/internal/domain"
)

func SendRequest(req *domain.Request, timeoutStr string, followRedirects bool, retryCount int) (*domain.Response, error) {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout}
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var resp *http.Response
	var respBody []byte

	for attempt := 0; attempt <= retryCount; attempt++ {
		httpReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
		if err != nil {
			return nil, err
		}
		httpReq.Header = req.Headers

		resp, err = client.Do(httpReq)
		if err != nil {
			color.Yellow("⚠️  Attempt %d failed: %v\n", attempt+1, err)
			if attempt == retryCount {
				return nil, err
			}
			time.Sleep(1 * time.Second) // optional backoff
			continue
		}

		respBody, _ = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		break // success
	}

	return &domain.Response{
		Status:  resp.Status,
		Headers: resp.Header,
		Body:    respBody,
	}, nil
}
