package handler

import (
	"net/http"
)

func (e *HttpEnv) TaskConfigHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, e.publicConfigFilePath)
}
