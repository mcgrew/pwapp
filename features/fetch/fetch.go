

//go:build fetch
package fetch

import (
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
//  	url := r.URL.Path[:11]
	response := []byte("Hello World!"
	w.Write(response)
}

func Post(w http.ResponseWriter, r *http.Request) {
	response := []byte("Hello World!"
	w.Write(response)
}

