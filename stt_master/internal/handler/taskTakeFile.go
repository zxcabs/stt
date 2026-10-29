package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"stt_master/internal/db"
)

const (
	TYPE_AUDIO = "audio"
	TYPE_SRT = "srt"
)

var validTypes = map[string]bool{
	TYPE_AUDIO: true,
	TYPE_SRT: true,
}

func (e *HttpEnv) TakeFileHandler(w http.ResponseWriter, r *http.Request) {
	if e.isServiceBusy() {
		e.sendJSONError(w, http.StatusServiceUnavailable, "Service is busy")
		return
	}

	id := r.URL.Query().Get("id")
	fileType := r.URL.Query().Get("type")

	if id == "" {
		e.sendJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	if !validTypes[fileType] {
		e.sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Wrong type '%s'", fileType))
		return
	}

	taskId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		e.sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Wrong id '%s'", id))
		return
	}

	task, tErr := db.GetTaskById(e.DB, taskId)

	if tErr != nil {
		e.sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Wrong id '%s'", id))
		return
	}

	switch fileType {
	case TYPE_AUDIO:
		e.sendAudioFile(task, w, r)
	default:
		e.sendJSONError(w, http.StatusNotImplemented, "This type not supported")
	}
}

func (e *HttpEnv) sendAudioFile(task *db.AudioTask, w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(e.audioDir, task.FilePath)
	http.ServeFile(w, r, filePath)
}
