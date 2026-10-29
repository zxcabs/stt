package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Init(dsn string) (*sql.DB, error) {

	if err := os.MkdirAll(filepath.Dir(dsn), os.ModePerm); err != nil {
		log.Printf("Error on make db dir %v\n", err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		fmt.Printf("DB Error %v\n", err)
		return nil, err
	}

	if err := createAudioTasksTable(db); err != nil {
		fmt.Printf("Ошибка создания таблицы: %v\n", err)
		defer db.Close()
		return nil, err
	}

	return db, nil
}
