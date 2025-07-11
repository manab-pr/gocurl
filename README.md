# 🚀 gocurl – A Simple Yet Powerful HTTP CLI Client in Go

[![Go](https://img.shields.io/badge/Go-1.21+-blue?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

`gocurl` is a modern command-line HTTP client written in Go — a fast, colorful, minimal alternative to tools like `curl` or `httpie`. Easily send HTTP requests, debug APIs, and pretty-print JSON responses with style.

---

## ✨ Features

| Feature                  | Description                                                                 |
|--------------------------|-----------------------------------------------------------------------------|
| `-X`                     | Specify request method: GET, POST, PUT, PATCH, DELETE                      |
| `-H`                     | Add custom headers (`-H "Auth: Bearer xyz"`)                                |
| `-d` / `@file.json`      | Add request body from string or file                                        |
| `-o filename.txt`        | Save response body to a file                                                |
| `--headers-only`         | Print only response headers                                                 |
| `--body-only`            | Print only response body                                                    |
| `--raw`                  | Skip JSON pretty-printing                                                   |
| `--follow`               | Follow redirects automatically                                              |
| `--timeout`              | Set custom timeout duration (e.g. `--timeout 5s`)                           |
| `--retry`                | Retry failed requests up to `n` times                                       |
| `--verbose`              | Print full request + response details                                       |
| ⏱ Response Timer         | Shows how long the request took with visual timing                         |
| 🌈 Colored Output        | Status, headers, and JSON output in color using `fatih/color`              |
| 📝 Logs                  | Saves request logs automatically                                            |

---

## 📦 Installation

### ✅ Via `go install`

```bash
go install github.com/manab-pr/gocurl/cmd/gocurl@latest


💻 Run Locally from Source

git clone https://github.com/manab-pr/gocurl.git
cd gocurl/cmd/gocurl
go run main.go https://jsonplaceholder.typicode.com/posts/1


🚦 Usage Examples

# Basic GET
gocurl https://jsonplaceholder.typicode.com/posts/1

# POST with data
gocurl -X POST -d '{"title":"Go is great!"}' https://jsonplaceholder.typicode.com/posts

# Use data from a file
gocurl -X POST -d @data.json https://jsonplaceholder.typicode.com/posts

# Add headers
gocurl -H "Authorization: Bearer token" https://example.com

# Save response
gocurl -o response.json https://jsonplaceholder.typicode.com/posts/1

# Follow redirects, set timeout, and retry
gocurl --follow --timeout 5s --retry 3 https://httpbin.org/redirect/2


🧪 Sample Output

🌍 Sending GET request to: https://jsonplaceholder.typicode.com/posts/1
✅ Status: 200 OK
📦 Headers:
   Content-Type: application/json; charset=utf-8
   ...

📄 Body:
{
  "userId": 1,
  "id": 1,
  "title": "Go is awesome!",
  "body": "This is a demo post"
}

⏳ Response time: 340.75ms █████


