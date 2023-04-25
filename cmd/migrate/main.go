package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	migrationsPath = "migrations"
	driver         = "postgres"
)

func main() {
	jww.SetLogThreshold(jww.LevelInfo)
	jww.SetStdoutThreshold(jww.LevelInfo)

	ctx := context.Background()
	jww.INFO.Println("Starting migrations")

	// Читает переменные окружения
	err := godotenv.Load()
	if err != nil {
		jww.INFO.Println("No .env file loaded")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = strings.ReplaceAll(os.Getenv("CI_PROJECT_NAME"), "-", "_")
	}
	if dbName == "" {
		dbName = "postgres"
	}
	dbDSN := fmt.Sprintf("host='%s' port=%s user='%s' password='%s' dbname='%s' sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		dbName)

	// подключаемся к БД
	db, err := connect(dbDSN)
	if err != nil {
		jww.ERROR.Fatalln(err)
	}
	jww.INFO.Println("The database connection was established successfully")

	// устанавливаем свой логер
	goose.SetLogger(&gooseLogger{ctx: ctx})
	_ = goose.SetDialect(driver)

	// запускаем миграции
	jww.INFO.Println("Upping migrations")
	err = goose.SetDialect("postgres")
	if err != nil {
		jww.ERROR.Fatalf("Failed to set dialect: %v", err)
	}
	err = goose.Up(db.DB(), migrationsPath)
	if err != nil {
		jww.ERROR.Fatalf("Failed to migrate: %v", err)
	}

	jww.INFO.Println("DB migration completed")
}

// Выполняет подключение к БД
func connect(dsn string) (*gorm.DB, error) {
	log.Println(dsn)
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Реализация интерфйса goose.Logger
type gooseLogger struct {
	ctx context.Context
}

func (gl *gooseLogger) Fatal(v ...interface{}) {
	jww.FATAL.Fatal(v...)
}
func (gl *gooseLogger) Fatalf(format string, v ...interface{}) {
	jww.FATAL.Fatalf(format, v...)
}
func (gl *gooseLogger) Print(v ...interface{}) {
	jww.INFO.Print(v...)
}
func (gl *gooseLogger) Println(v ...interface{}) {
	jww.INFO.Println(v...)
}
func (gl *gooseLogger) Printf(format string, v ...interface{}) {
	jww.INFO.Printf(format, v...)
}
