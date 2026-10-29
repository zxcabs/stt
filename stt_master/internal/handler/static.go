package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type staticHandler struct {
	dirPath   string
	indexHame string
}

func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, filepath.Join(h.dirPath, h.indexHame))
		return
	}

	path := filepath.Join(h.dirPath, r.URL.Path)

	fi, err := os.Stat(path)

	if os.IsNotExist(err) || fi.IsDir() {
		fmt.Printf("Not found %v\n", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Printf("Internal error %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.dirPath)).ServeHTTP(w, r)
}
