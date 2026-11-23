
//go:build !files
package files

import (
	"net/http"
)

const Enabled = false

func Get(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func Put(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func Post(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func Delete(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
