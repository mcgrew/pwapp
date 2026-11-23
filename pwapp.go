package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"pwapp/features/fetch"
	"pwapp/features/files"
)

//go:embed html
var webContent embed.FS

//go:embed pwapp.js
var pwapplib []byte

// command line flags
var headless bool
var hostPort int
var hostIp = "[::1]"
var workDir string

const pingScript = "setInterval(function() {fetch('_pwping').then()}, 7000);"

func Ternary[T any](cond bool, yes, no T) T {
	if cond { return yes }; return no
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// only allow GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	path := "html" + r.URL.Path
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	log.Printf("Fetching %s\n", path)
	data, err := webContent.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("Not Found: %s\n", path)
	}
	mimetype := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mimetype)
	w.Write(injectJS(path, data))
}

func injectJS(path string, html []byte) []byte {
	if headless || !strings.HasSuffix(path, "html"){ return html }
	injection := []byte( `</body><script type="text/javascript">` +
		pingScript + fetch.FetchScript + "</script>")
	return bytes.Replace(html, []byte("</body>"), injection, 1)
}

func launchBrowser() error {
	var cmd string
	var args []string
	url := fmt.Sprintf("http://%s:%d/", hostIp, hostPort)

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

func parseFlags() {
	var ipv4 bool
	flag.IntVar(&hostPort, "port", 37777, "The port number to run on")
	flag.BoolVar(&headless, "headless", false, "Do not launch the web browser")
	flag.BoolVar(&ipv4, "ipv4", false, "Use ipv4 instead of ipv6")
	flag.StringVar(&workDir, "workdir", "", "The working directory for the server")
	flag.Parse()
	if ipv4 { hostIp = "127.0.0.1" }
}

func startUI(srv *http.Server, mux *http.ServeMux) {
	if headless {
		return
	}
	// listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// set up a ping endpoint unless we are in headless mode.
	// this will shut down the server if the browser tab is closed.
	pingTimeout := 10 * time.Second
	killTimer := time.NewTimer(pingTimeout)
	defer killTimer.Stop()
	mux.HandleFunc("GET /_pwping", func(w http.ResponseWriter, r *http.Request) {
		killTimer.Reset(pingTimeout)
	})
	launchBrowser()

	select {
		case <-killTimer.C:
			log.Printf("No ping received. Shutting down...")
		case <-ctx.Done():
			log.Println("OS interrupt received. Shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func main() {
	parseFlags()
	if workDir != "" {
		err := os.Chdir(workDir)
		if err != nil {
			log.Fatalf("Unable to open directory %s", workDir)
			log.Fatal(err)
			os.Exit(-1)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	if fetch.Enabled {
		mux.HandleFunc("/_pwfetch/", fetch.Request)
	}
	if files.Enabled {
		mux.HandleFunc("GET /_pwfile/", files.Get) // read file/dir
		mux.HandleFunc("PUT /_pwfile/", files.Put) // write file/dir
		mux.HandleFunc("POST /_pwfile/", files.Post) // move file
		mux.HandleFunc("DELETE /_pwfile/", files.Delete) // delete file
	}
	mux.HandleFunc("GET /_pwapp.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(pwapplib)
	})


	bindInterface := fmt.Sprintf("%s:%d", Ternary(headless, "0.0.0.0", hostIp), hostPort)
	// Start the server
	srv := &http.Server{
		Addr:		  bindInterface,
		Handler:	  mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go startUI(srv, mux)

	log.Printf("Server listening on %s", bindInterface)
	srv.ListenAndServe()
}

