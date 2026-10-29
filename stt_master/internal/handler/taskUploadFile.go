package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"stt_master/internal/db"
)

const maxUploadFileSize = 100 * 1024 * 1024

func (e *HttpEnv) TaskUploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadFileSize)

	if r.Method != http.MethodPost {
		e.sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	task, err := db.GetTaskById(e.DB, taskId)

	if err != nil {
		e.sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Wrong id '%s'", id))
		return
	}

	if task.Status != db.AUDIO_TASKS_SATUS_PROCESSING {
		e.sendJSONError(w, http.StatusBadRequest, "Task is not in processing")
		return
	}

	switch fileType {
	case TYPE_SRT:
		e.uploadSrtFile(task, w, r)
	default:
		e.sendJSONError(w, http.StatusNotImplemented, "This type not supported")
	}
}

func (e *HttpEnv) uploadSrtFile(task *db.AudioTask, w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")

	if err != nil {
		e.sendJSONError(w, http.StatusBadRequest, "Error on read file")
		return
	}

	defer file.Close()

	tempFileName := rand.Text() + ".srt"

	if err := os.MkdirAll(e.tempDir, os.ModePerm); err != nil {
		e.sendJSONError(w, http.StatusInternalServerError, "Error")
		return
	}

	tempFilePath := filepath.Join(e.tempDir, tempFileName)
	tempFile, err := os.Create(tempFilePath)

	if err != nil {
		e.sendJSONError(w, http.StatusInternalServerError, "Error")
		return
	}

	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		e.sendJSONError(w, http.StatusInternalServerError, "Error")
		return
	}
	tempFile.Close()
	file.Close()

	destFilePath := filepath.Join(e.srtDir, changeExt(task.FilePath, ".srt"))

	if err := moveFile(tempFilePath, destFilePath); err != nil {
		fmt.Printf("Error on move file %v\n", err)
		e.sendJSONError(w, http.StatusInternalServerError, "Error")
	}

	if err := db.MarkCompletedTaskById(e.DB, task.ID); err != nil {
		e.sendJSONError(w, http.StatusInternalServerError, "Error")
		return
	}

	e.sendJSON(w, http.StatusOK, "ok")
}

func moveFile(oldPath, newPath string) error {
	destDir := filepath.Dir(newPath)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.Rename(oldPath, newPath); err == nil {
		return nil
	}

	src, err := os.Open(oldPath)
	if err != nil {
		return fmt.Errorf("Error open src file %w", err)
	}
	defer src.Close()

	dest, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("Error open dest file %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("Error on copy file %w", err)
	}

	src.Close()
	dest.Close()

	if err := os.Remove(oldPath); err != nil {
		fmt.Printf("Warning error on remove temp file %v", err)
	}

	return nil
}
