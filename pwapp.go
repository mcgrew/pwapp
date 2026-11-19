package main

import (
	"embed"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//go:embed html
var CONTENT embed.FS

var PORT = 7777

type HttpHandler struct {
	Writer http.ResponseWriter
	request *http.Request
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := "html" + r.URL.Path
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	data, err := CONTENT.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		fmt.Printf("Not Found: %s\n", path)
	}
	fmt.Printf("Fetching %s\n", path)
	mimetype := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mimetype)
	w.Write(data)
}

func external(url string) string {
	return ""
}


func launchBrowser() error {
	var cmd string
	var args []string
	url := fmt.Sprintf("http://127.0.0.1:%d/", PORT)

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}


func main() {

	handler := HttpHandler{}
	// Start the server on port 8080
	srv := &http.Server{
		Addr:		 fmt.Sprintf(":%d", PORT),
		Handler:	  handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	launchBrowser()
	log.Println("Server listening on :8080")
	log.Fatal(srv.ListenAndServe())
}

