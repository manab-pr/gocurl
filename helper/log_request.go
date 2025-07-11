package helper

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func LogRequest(method, url, body string, headers http.Header) {
	f, err := os.OpenFile(".gocurl-history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("❌ Failed to log request: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var logData strings.Builder
	logData.WriteString(fmt.Sprintf("[%s] %s %s\n", timestamp, method, url))

	if len(headers) > 0 {
		logData.WriteString("Headers:\n")
		for k, v := range headers {
			logData.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
		}
	}

	if strings.TrimSpace(body) != "" {
		logData.WriteString(fmt.Sprintf("Body:\n  %s\n", body))
	}

	logData.WriteString("---\n\n")

	if _, err := f.WriteString(logData.String()); err != nil {
		color.Red("❌ Failed to write request log: %v\n", err)
	}
}
