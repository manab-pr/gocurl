package application

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/manab-pr/gocurl/internal/domain"
	"github.com/manab-pr/gocurl/internal/infrastructure/cli"
	infraHttp "github.com/manab-pr/gocurl/internal/infrastructure/http"
	"github.com/manab-pr/gocurl/internal/infrastructure/logger"
)

func Run() {
	cli.ParseFlags()

	if !cli.NoBanner {
		cli.ShowBanner()
	}

	if len(os.Args) < 2 {
		fmt.Println("🚨 Usage: gocurl [options] URL")
		return
	}

	url := os.Args[len(os.Args)-1]

	req := domain.NewRequestBuilder().
		WithMethod(cli.Method).
		WithURL(url).
		WithBody(cli.Data).
		WithHeaders(cli.Headers).
		Build()

	if cli.ExportCurl {
		cli.PrintCurlCommand(req.Method, req.URL, req.Body, req.Headers)
	}

	if cli.Verbose {
		cli.PrintVerbose(req)
	}

	logger.LogRequest(req.Method, req.URL, req.Body, req.Headers)

	resp, err := infraHttp.SendRequest(req, cli.TimeoutStr, cli.FollowRedirects, cli.RetryCount)
	if err != nil {
		color.Red("❌ Request failed: %v\n", err)
		return
	}

	cli.PrintResponse(resp, cli.RawOutput, cli.HeadersOnly, cli.BodyOnly, cli.OutFile)
}
