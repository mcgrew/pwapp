
//go:build !fetch
package fetch

import (
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Post(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
