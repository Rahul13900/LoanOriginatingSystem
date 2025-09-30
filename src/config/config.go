package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DBURL       string
	WorkerCount int
}

func LoadFromEnv() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db := os.Getenv("DB_URL")
	if db == "" {
		db = ""
	}
	workerCount := 0
	if wc := os.Getenv("WORKER_COUNT"); wc != "" {
		if n, err := strconv.Atoi(wc); err == nil {
			workerCount = n
		}
	} else {
		workerCount = 4
	}

	return &Config{
		Port:        port,
		DBURL:       db,
		WorkerCount: workerCount,
	}
}
