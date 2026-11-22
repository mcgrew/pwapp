

//go:build fetch
package fetch

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const Enabled = true

func Request(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.Path[11:]
	if !strings.HasPrefix(reqUrl, "http") {
		http.NotFound(w, r)
		return
	}

	// construct a new request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	proxyReq, err := http.NewRequestWithContext(ctx, r.Method, reqUrl, nil)
	if err != nil {
		log.Fatalf("Failed to build request: %v", err)
	}
	if err != nil {
		log.Printf("Could not create request for %s", reqUrl)
		http.NotFound(w, r)
		return
	}
	// copy over the request headers
	for key, values := range r.Header {
		for _, v := range values {
			proxyReq.Header.Add(key, v)
		}
	}

	// send the request to the intended recipient
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error requesting %s", reqUrl)
		http.NotFound(w, r)
		return
	}

	// send the response back to the original requester
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	log.Println(resp.Body)
	io.Copy(w, resp.Body)
}

