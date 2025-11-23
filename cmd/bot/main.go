package main

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/LigeronAhill/domovenok/internal/config"
	"github.com/LigeronAhill/domovenok/internal/logger"
)

func main() {
	env := "local"
	level := logger.DEBUG
	if prod := os.Getenv("ENVIRONMENT"); prod != "" {
		if strings.ToLower(prod) == "production" {
			env = "production"
			level = logger.INFO
		}
	}
	logger.Init(level)
	settings, err := config.Init(env)
	if err != nil {
		log.Fatal(err)
	}
	if settings == nil {
		log.Fatal("не удалось загрузить конфигурацию")
	}
	slog.Info("Конфигурация успешно загружена", slog.String("Окружение", env))
}
