// main.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("🚨 Usage: gocurl [get|post] URL")
		return
	}

	method := strings.ToUpper(os.Args[1])
	url := os.Args[2]

	start := time.Now()

	switch method {
	case "GET":
		doGet(url)
	case "POST":
		doPost(url)
	default:
		fmt.Println("❌ Only GET and POST supported for now")
	}

	elapsed := time.Since(start)
	fmt.Printf("⏱️  Done in %s\n", elapsed)
}

func doGet(url string) {
	fmt.Println("🌍 Sending GET request to:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("🚨 Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅ Status:", resp.Status)
	fmt.Println("📦 Headers:")
	for k, v := range resp.Header {
		fmt.Printf("   %s: %s\n", k, strings.Join(v, ", "))
	}

	fmt.Println("\n📄 Body:")
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func doPost(url string) {
	fmt.Println("🌍 Sending POST request to:", url)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("🚨 Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅ Status:", resp.Status)
	fmt.Println("\n📄 Body:")
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
