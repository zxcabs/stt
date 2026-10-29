package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stt_master/internal/db"
	"stt_master/internal/handler"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	SERVER_HOST string
	SERVER_PORT string
	DB_DSN      string
}

func getEnvDef(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not present")
	}

	conf := Config{
		SERVER_HOST: getEnvDef("SERVER_HOST", ""),
		SERVER_PORT: getEnvDef("SERVER_PORT", "8080"),
		DB_DSN:      getEnvDef("DB_DSN", "./db/example.db?_journal_mode=WAL&_busy_timeout=5000"),
	}

	db, errDb := db.Init(conf.DB_DSN)

	if errDb != nil {
		log.Printf("Критическая ошибка бд: %v", errDb)
		os.Exit(1)
	}

	httpEnv := handler.New(db)
	router := httpEnv.RegisterRoutes()

	addr := fmt.Sprintf("%s:%s", conf.SERVER_HOST, conf.SERVER_PORT)

	log.Printf("Server starting on %q", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatalf("Error %v", err)
	}
}
