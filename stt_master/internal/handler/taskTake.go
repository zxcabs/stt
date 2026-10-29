package handler

import (
	"net/http"
	"stt_master/internal/db"
)

func (e *HttpEnv) TakeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if e.isServiceBusy() {
		e.sendJSONError(w, http.StatusServiceUnavailable, "Service is busy")
		return
	}

	task, err := db.ClaimNextTask(e.DB)

	if err != nil {
		e.sendJSONError(w, http.StatusBadGateway, "Some error")
		return
	}

	e.sendJSON(w, http.StatusOK, task)
}
