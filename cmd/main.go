package main

import (
	"errors"
	"fmt"
	"go-final-project/internal/config"
	"go-final-project/internal/db"
	"go-final-project/internal/server"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("БД запущена. Путь: %s\n", cfg.DBFile)
	defer func() {
		if db.DB != nil {
			_ = db.DB.Close()
		}
	}()

	srvCfg := server.DefaultConfig()
	srvCfg.Port = cfg.Port

	fmt.Printf("Запускаем сервер на порту %d\n", srvCfg.Port)
	if err := server.NewServer(srvCfg).ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
