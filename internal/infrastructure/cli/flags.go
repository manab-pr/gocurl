package cli

import (
	"flag"
	"fmt"
)

var (
	Method          string
	Headers         HeaderFlags
	Data            string
	OutFile         string
	RawOutput       bool
	HeadersOnly     bool
	BodyOnly        bool
	Verbose         bool
	FollowRedirects bool
	TimeoutStr      string
	RetryCount      int
	ExportCurl      bool
	NoBanner        bool
)

func ParseFlags() {
	flag.StringVar(&Method, "X", "GET", "HTTP method to use")
	flag.Var(&Headers, "H", "Custom header (can be used multiple times)")
	flag.StringVar(&Data, "d", "", "Request body data (for POST, PUT, PATCH)")
	flag.StringVar(&OutFile, "o", "", "Output file to save response body")
	flag.BoolVar(&RawOutput, "raw", false, "Show raw output even if JSON")
	flag.BoolVar(&HeadersOnly, "headers-only", false, "Show only response headers")
	flag.BoolVar(&BodyOnly, "body-only", false, "Show only response body")
	flag.BoolVar(&Verbose, "verbose", false, "Print full request details (method, headers, body)")
	flag.BoolVar(&Verbose, "v", false, "Shorthand for --verbose")
	flag.BoolVar(&FollowRedirects, "follow", false, "Follow HTTP redirects")
	flag.StringVar(&TimeoutStr, "timeout", "10s", "Set request timeout (e.g. 5s, 2m)")
	flag.IntVar(&RetryCount, "retry", 0, "Retry failed requests up to n times")
	flag.BoolVar(&ExportCurl, "export-curl", false, "Print the equivalent curl command")
	flag.BoolVar(&NoBanner, "no-banner", false, "Disable ASCII banner")
	flag.Parse()
}

type HeaderFlags []string

func (h *HeaderFlags) String() string {
	return fmt.Sprintf("%v", *h)
}
func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}
