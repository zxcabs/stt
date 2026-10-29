package handler

import (
	"net/http"
	"stt_master/internal/db"
)

func (e *HttpEnv) StatusHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetAllTasks(e.DB)

	if err != nil {
		e.sendJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	e.sendJSON(w, http.StatusOK, tasks)
}
