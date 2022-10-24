package main

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/pkg/app"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.WithError(err).Error("cant create new cfg")

		os.Exit(2)
	}

	ctx = config.WrapContext(ctx, cfg) // обогащаем конфиг контекстом
	// Создание приложения
	application, err := app.NewApp(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("app creating error")

		os.Exit(2)
	}
	// Запуск приложения
	err = application.Run(ctx)
	if err != nil {
		log.WithError(err).Error("app startup error")

		os.Exit(2)
	}
}
