package config

import (
	"os"
	"strconv"
)

const (
	defaultPort   = 7540
	defaultDBFile = "scheduler.db"
)

type Config struct {
	Port   int
	DBFile string
}

func Load() Config {
	cfg := Config{
		Port:   defaultPort,
		DBFile: defaultDBFile,
	}

	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		cfg.DBFile = envDBFile
	}

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if parsedPort, err := strconv.Atoi(envPort); err == nil && parsedPort > 0 {
			cfg.Port = parsedPort
		}
	}

	return cfg
}
