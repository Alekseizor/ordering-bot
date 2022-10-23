package main

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/pkg/app"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	SetFormatter(&JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	SetOutput(os.Stdout)

	// Only log the warning severity or above.
	SetLevel(DebugLevel)
}

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		WithError(err).Error("cant create new cfg")

		os.Exit(2)
	}

	ctx = config.WrapContext(ctx, cfg) // обогащаем конфиг контекстом
	// Создание приложения
	application, err := app.NewApp(ctx)
	if err != nil {
		WithContext(ctx).WithError(err).Error("app creating error")

		os.Exit(2)
	}
	// Запуск приложения
	err = application.Run(ctx)
	if err != nil {
		WithError(err).Error("app startup error")

		os.Exit(2)
	}
}
