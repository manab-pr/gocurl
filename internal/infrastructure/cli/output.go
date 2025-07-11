package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manab-pr/gocurl/internal/domain"
	"github.com/manab-pr/gocurl/pkg/utils"
)

func PrintResponse(resp *domain.Response, rawOutput, headersOnly, bodyOnly bool, outFile string) {
	if headersOnly && bodyOnly {
		color.Red("❌ You cannot use --headers-only and --body-only together.")
		return
	}

	if !bodyOnly {
		color.Green("✅ Status: %s\n", resp.Status)
		color.Yellow("📦 Headers:")
		for k, v := range resp.Headers {
			fmt.Printf("   %s: %s\n", k, strings.Join(v, ", "))
		}
	}

	if !headersOnly {
		fmt.Println("\n📄 Body:")
		if outFile != "" {
			os.WriteFile(outFile, resp.Body, 0644)
			color.Blue("💾 Response saved to %s\n", outFile)
		} else if strings.Contains(resp.Headers.Get("Content-Type"), "application/json") {
			if rawOutput {
				fmt.Println(string(resp.Body))
			} else {
				prettyPrintJSON(resp.Body)
			}
		} else {
			fmt.Println(string(resp.Body))
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

	utils.PrintColoredJSON(raw, 0)
}

func PrintCurlCommand(method, url, body string, headers http.Header) {
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

func PrintVerbose(req *domain.Request) {
	color.Blue("\n🔍 Verbose Request Info:")
	fmt.Printf("➡️  Method: %s\n", req.Method)
	fmt.Printf("🔗 URL: %s\n", req.URL)
	color.Blue("🧾 Headers:")
	for k, v := range req.Headers {
		fmt.Printf("   %s: %s\n", k, strings.Join(v, ", "))
	}
	if req.Body != "" {
		color.Blue("📤 Body:\n%s\n", req.Body)
	}
}