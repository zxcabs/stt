package handler

import (
	"database/sql"
	"sync"
)

type HttpEnv struct {
	DB                   *sql.DB
	tempDir              string
	audioDir             string
	srtDir               string
	publicConfigFilePath string
	isIndexing           bool
	indexingMu           sync.RWMutex
}

func New(db *sql.DB) *HttpEnv {
	return &HttpEnv{
		DB:                   db,
		tempDir:              "/data/temp",
		audioDir:             "/data/original_audio",
		srtDir:               "/data/original_srt",
		publicConfigFilePath: "/data/cfg.json",
	}
}

func (e *HttpEnv) isServiceBusy() bool {
	e.indexingMu.Lock()

	if e.isIndexing {
		e.indexingMu.Unlock()
		return true
	}

	e.indexingMu.Unlock()

	return false
}
