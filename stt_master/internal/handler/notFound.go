package handler

import (
	"net/http"
)

func (e *HttpEnv) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}
