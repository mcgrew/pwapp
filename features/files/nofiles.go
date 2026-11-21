
//go:build !files
package files

import (
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Write(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Move(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func MkDir(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
