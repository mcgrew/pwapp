
//go:build files
package files

import (
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[6:]
	fmt.Printf("Write File %s\n", path)
	fmt.Fprintf(w, "Hello from the file read api!")
}

func Write(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[6:]
	fmt.Printf("Write File %s\n", path)
	fmt.Fprintf(w, "Hello from the file write api!")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[6:]
	fmt.Printf("Write File %s\n", path)
	fmt.Fprintf(w, "Hello from the file delete api!")
}

func Move(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[6:]
	fmt.Printf("Write File %s\n", path)
	fmt.Fprintf(w, "Hello from the file move api!")
}

func MkDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[6:]
	fmt.Printf("Write File %s\n", path)
	fmt.Fprintf(w, "Hello from the directory create api!")
}
