
//go:build files
package files

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const Enabled = true
const pathStart = 9

// Handler functions
func Get(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) { return }
	readFile(w, r)
}

func Put(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) { return }
	path := r.URL.Path[pathStart:]
	if strings.HasSuffix(path, "/") {
		mkDir(w, r)
		return
	}
	writeFile(w, r)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) { return }
	deleteFile(w, r)
}

func Post(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) { return }
	postMove(w, r)
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
//  	path := r.URL.Path[pathStart:]
//  	if false {
//  		http.Error(w, "Forbidden", http.StatusForbidden)
//  		return false
//  	}
	return true
}

// Operation functions
func readFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	if strings.HasSuffix(path, "/") {
		listDir(w, r)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("Read file %s", path)
	mimetype := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mimetype)
	w.Write(data)
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	if strings.HasSuffix(path, "/") {
		mkDir(w, r)
		return
	}
	log.Printf("Writing %s", path)
	body, _ := io.ReadAll(r.Body)
	err := os.WriteFile(path, body, 0644)
	if err != nil {
		message := fmt.Sprintf("Unable to write file %s", err)
		log.Fatal(message)
		log.Fatal(err)
		http.Error(w, message, http.StatusBadRequest)
	}
	w.Header().Set("Location", r.URL.Path)
	w.WriteHeader(http.StatusCreated)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	log.Printf("Deleting %s", path)
	err := os.Remove(path)
	if err != nil {
		message := fmt.Sprintf("Unable to delete file %s", path)
		log.Fatal(message)
		log.Fatal(err)
		http.Error(w, message, http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
}

func postMove(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	body, _ := io.ReadAll(r.Body)
	dest := string(body)
	log.Printf("Move %s -> %s", path, dest)
	err := os.Rename(path, dest)
	if err != nil {
		message := fmt.Sprintf("Unable to move file %s to %s", path, dest)
		log.Fatal(message)
		log.Fatal(err)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", r.URL.Path[:9] + dest)
	w.WriteHeader(http.StatusCreated)
}

func mkDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	err := os.MkdirAll(path, 0755)
	if err != nil {
		message := fmt.Sprintf("Unable to create directory %s", path)
		log.Fatal(message)
		log.Fatal(err)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", r.URL.Path)
	w.WriteHeader(http.StatusCreated)
}

func listDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[pathStart:]
	entries, err := os.ReadDir(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	names := `["`
	for _, entry := range entries {
		if len(names) > 2 { names += `","` }
		names += entry.Name()
	}
	names = names[:len(names)-1] + `"]`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(names))
}
