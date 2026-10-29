package handler

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"stt_master/internal/db"

	"golang.org/x/sync/errgroup"
)

func (e *HttpEnv) TasksToIndexHandler(w http.ResponseWriter, r *http.Request) {
	e.indexingMu.Lock()

	if e.isIndexing {
		e.indexingMu.Unlock()
		e.sendJSONError(w, http.StatusConflict, "Tasks to index already runnig")

		return
	}

	e.isIndexing = true
	e.indexingMu.Unlock()

	go e.runDiskIndexing()

	e.sendJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Starting indexing",
	})
}

func (e *HttpEnv) runDiskIndexing() error {
	defer func() {
		e.indexingMu.Lock()
		e.isIndexing = false
		e.indexingMu.Unlock()
		log.Println("Done indexing audio tasks")
	}()

	log.Println("Starting indexing audio tasks")

	if err := db.RecreateAudioTasksTable(e.DB); err != nil {
		return fmt.Errorf("Error on create table %w", err)
	}

	var audioFiles, srtFiles map[string]bool
	var g errgroup.Group

	g.Go(func() error {
		var err error
		audioFiles, err = scanDir(e.audioDir, ".wav")
		return err
	})

	g.Go(func() error {
		var err error
		srtFiles, err = scanDir(e.srtDir, ".srt")
		return err
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("Error on indexing %w", err)
	}

	for audioPath := range audioFiles {
		srtPath := changeExt(audioPath, ".srt")
		status := db.AUDIO_TASKS_SATUS_PENDING

		if srtFiles[srtPath] == true {
			status = db.AUDIO_TASKS_SATUS_COMPLETED
		}

		if err := db.AddTask(e.DB, audioPath, status); err != nil {
			return nil
		}
	}

	return nil
}

func scanDir(root string, ext string) (map[string]bool, error) {
	if root == "" || ext == "" {
		return nil, fmt.Errorf("Rootpath '%s' or ext '%s' is empty", root, ext)
	}

	absRoot, err := filepath.Abs(root)

	if err != nil {
		return nil, err
	}

	files := make(map[string]bool)
	normalizedExt := strings.ToLower(ext)

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.IsDir() && strings.ToLower(filepath.Ext(d.Name())) == normalizedExt {
			relPath, err := filepath.Rel(absRoot, path)

			if err != nil {
				return err
			}

			files[relPath] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func changeExt(name, newExt string) string {
	oldExt := filepath.Ext(name)
	base := strings.TrimSuffix(name, oldExt)

	if newExt != "" && !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}

	return base + newExt
}
