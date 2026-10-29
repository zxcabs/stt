package handler

import (
	"fmt"
	"net/http"
)

func (e *HttpEnv) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}
