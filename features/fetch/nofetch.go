
//go:build !fetch
package fetch

import (
	"net/http"
)

const Enabled = false

func Request(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
