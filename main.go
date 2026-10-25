// Minimal HTTP server that reports the CPU architecture it runs on.
// Used to confirm that one multi-architecture tag runs natively on
// both amd64 and arm64 hosts.
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"runtime"
)

const page = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>multiarch-demo</title>
  <style>
    body { font-family: system-ui, sans-serif; display: grid; place-items: center; min-height: 100vh; margin: 0; background: #f4f4f7; color: #232f3e; }
    main { background: #fff; padding: 2rem 3rem; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.08); text-align: center; }
    .arch { font-size: 3rem; font-weight: 700; color: #4d27b2; margin: .5rem 0; }
    p { margin: .25rem 0; }
  </style>
</head>
<body>
  <main>
    <p>This container is running on</p>
    <p class="arch">%s</p>
    <p>OS: %s</p>
    <p>Host: %s</p>
  </main>
</body>
</html>
`

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	mux := http.NewServeMux()

	// Plain-text architecture, used by the verify job in the workflow.
	mux.HandleFunc("/arch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, runtime.GOARCH)
	})

	// Health check endpoint for container orchestrators.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, page, runtime.GOARCH, runtime.GOOS, html.EscapeString(hostname))
	})

	addr := ":8080"
	log.Printf("listening on %s (%s/%s)", addr, runtime.GOOS, runtime.GOARCH)
	log.Fatal(http.ListenAndServe(addr, mux))
}
