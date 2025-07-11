package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manab-pr/gocurl/helper"
)

var (
	method          string
	headers         headerFlags
	data            string
	outFile         string
	rawOutput       bool
	headersOnly     bool
	bodyOnly        bool
	verbose         bool
	followRedirects bool
	timeoutStr      string
	retryCount      int
	exportCurl      bool
	noBanner        bool
)

func init() {
	flag.StringVar(&method, "X", "GET", "HTTP method to use")
	flag.Var(&headers, "H", "Custom header (can be used multiple times)")
	flag.StringVar(&data, "d", "", "Request body data (for POST, PUT, PATCH)")
	flag.StringVar(&outFile, "o", "", "Output file to save response body")
	flag.BoolVar(&rawOutput, "raw", false, "Show raw output even if JSON")
	flag.BoolVar(&headersOnly, "headers-only", false, "Show only response headers")
	flag.BoolVar(&bodyOnly, "body-only", false, "Show only response body")
	flag.BoolVar(&verbose, "verbose", false, "Print full request details (method, headers, body)")
	flag.BoolVar(&verbose, "v", false, "Shorthand for --verbose")
	flag.BoolVar(&followRedirects, "follow", false, "Follow HTTP redirects")
	flag.StringVar(&timeoutStr, "timeout", "10s", "Set request timeout (e.g. 5s, 2m)")
	flag.IntVar(&retryCount, "retry", 0, "Retry failed requests up to n times")
	flag.BoolVar(&exportCurl, "export-curl", false, "Print the equivalent curl command")
	flag.BoolVar(&noBanner, "no-banner", false, "Disable ASCII banner")

}

func main() {
	flag.Parse()

	if !noBanner {
		showBanner()
	}

	if flag.NArg() < 1 {
		fmt.Println("🚨 Usage: gocurl [options] URL")
		return
	}

	if flag.NArg() < 1 {
		fmt.Println("🚨 Usage: gocurl [options] URL")
		return
	}

	url := flag.Arg(0)

	method = strings.ToUpper(method)
	start := time.Now()
	sendRequest(method, url, data)
	elapsed := time.Since(start)
	color.New(color.FgHiBlack).Printf("⏱️  Done in %s\n", elapsed)
}

func sendRequest(method, url, body string) {
	color.Cyan("🌍 Sending %s request to: %s\n", method, url)

	// Load file if -d @file.json
	if strings.HasPrefix(body, "@") {
		filename := strings.TrimPrefix(body, "@")
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			color.Red("❌ Failed to read file %s: %v\n", filename, err)
			return
		}
		body = string(fileContent)
	}

	// Parse timeout
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		color.Red("❌ Invalid timeout duration: %v\n", err)
		return
	}

	client := &http.Client{Timeout: timeout}
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Retry loop
	var resp *http.Response
	var respBody []byte

	for attempt := 0; attempt <= retryCount; attempt++ {
		req, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			color.Red("❌ Request creation failed: %v\n", err)
			return
		}

		// Set headers
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				req.Header.Add(key, val)
			}
		}

		if body != "" && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		if exportCurl {
			printCurlCommand(method, url, body, req.Header)
		}

		// Verbose print
		if verbose {
			color.Blue("\n🔍 Verbose Request Info:")
			fmt.Printf("➡️  Method: %s\n", req.Method)
			fmt.Printf("🔗 URL: %s\n", req.URL.String())
			color.Blue("🧾 Headers:")
			for k, v := range req.Header {
				fmt.Printf("   %s: %s\n", k, strings.Join(v, ", "))
			}
			if body != "" {
				color.Blue("📤 Body:\n%s\n", body)
			}
		}

		// 🔁 Make request
		resp, err = client.Do(req)
		originalBody := body // store it before changing
		helper.LogRequest(method, url, originalBody, req.Header)

		if err != nil {
			color.Yellow("⚠️  Attempt %d failed: %v\n", attempt+1, err)
			if attempt == retryCount {
				color.Red("❌ All retry attempts failed.\n")
				return
			}
			time.Sleep(1 * time.Second) // optional backoff
			continue
		}

		respBody, _ = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		break // success
	}

	// Output
	if headersOnly && bodyOnly {
		color.Red("❌ You cannot use --headers-only and --body-only together.")
		return
	}

	if !bodyOnly {
		color.Green("✅ Status: %s\n", resp.Status)
		color.Yellow("📦 Headers:")
		for k, v := range resp.Header {
			fmt.Printf("   %s: %s\n", k, strings.Join(v, ", "))
		}
	}

	if !headersOnly {
		fmt.Println("\n📄 Body:")
		if outFile != "" {
			os.WriteFile(outFile, respBody, 0644)
			color.Blue("💾 Response saved to %s\n", outFile)
		} else if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			if rawOutput {
				fmt.Println(string(respBody))
			} else {
				prettyPrintJSON(respBody)
			}
		} else {
			fmt.Println(string(respBody))
		}
	}
}

func prettyPrintJSON(data []byte) {
	var raw interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		fmt.Println(string(data)) // fallback
		return
	}

	helper.PrintColoredJSON(raw, 0)
}

type headerFlags []string

func (h *headerFlags) String() string {
	return fmt.Sprintf("%v", *h)
}
func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func printTimeBar(duration time.Duration) {
	ms := duration.Milliseconds()

	var blocks int
	switch {
	case ms < 100:
		blocks = 1
	case ms < 300:
		blocks = 2
	case ms < 600:
		blocks = 4
	case ms < 1000:
		blocks = 6
	case ms < 2000:
		blocks = 8
	default:
		blocks = 10
	}

	bar := strings.Repeat("█", blocks)
	color.Cyan("⏳ Response time: %s %s\n", duration, bar)
}

func printCurlCommand(method, url, body string, headers http.Header) {
	var b strings.Builder
	b.WriteString("curl")

	// Method
	if method != "GET" {
		b.WriteString(fmt.Sprintf(" -X %s", method))
	}

	// Headers
	for k, v := range headers {
		for _, val := range v {
			b.WriteString(fmt.Sprintf(` -H "%s: %s"`, k, val))
		}
	}

	// Body
	if strings.TrimSpace(body) != "" {
		b.WriteString(fmt.Sprintf(` -d '%s'`, body))
	}

	// URL
	b.WriteString(" " + url)

	color.Cyan("\n📤 Equivalent curl command:\n%s\n", b.String())
}

func showBanner() {
	banner := `
   __  __    _    _   _    _    ____  
  |  \/  |  / \  | \ | |  / \  | __ ) 
  | |\/| | / _ \ |  \| | / _ \ |  _ \ 
  | |  | |/ ___ \| |\  |/ ___ \| |_) |
  |_|  |_/_/   \_\_| \_/_/   \_\____/ 

`
	color.HiMagenta(banner)
}
