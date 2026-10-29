package handler

import (
	"net/http"
)

func (e *HttpEnv) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", staticHandler{dirPath: "./web", indexHame: "index.html"})
	mux.HandleFunc("GET /api/tasks/status", e.StatusHandler)
	mux.HandleFunc("GET /api/tasks/config", e.TaskConfigHandler)
	mux.HandleFunc("POST /api/tasks/take", e.TakeTaskHandler)
	mux.HandleFunc("GET /api/tasks/take/file", e.TakeFileHandler)
	mux.HandleFunc("POST /api/tasks/upload/file", e.TaskUploadFileHandler)
	mux.HandleFunc("POST /api/tasks/indexing", e.TasksToIndexHandler)

	mux.HandleFunc("/", e.NotFoundHandler)

	handlerWithLog := LoggerMiddleware(mux)

	return handlerWithLog
}
