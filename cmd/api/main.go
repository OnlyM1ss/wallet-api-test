package main

import (
	"log"
	"walletapitest/internal/app"
	"walletapitest/internal/config"
	"walletapitest/internal/pkg/logger"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logger := logger.New(cfg.LogLevel)

	// Создание и запуск приложения
	application := app.New(cfg, logger)
	
	if err := application.Run(); err != nil {
		logger.Fatal("Application failed to start", "error", err)
	}
}