package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Btlachac/discord-trivia-bot/postgres"
	"github.com/Btlachac/discord-trivia-bot/service"
	"github.com/Btlachac/discord-trivia-bot/worker/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	db := getDBConnection()
	runMigrations(db)

	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	logger, err := newZapLogger(true)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create zap logger > %w", err))
	}
	ts := createTriviaService(db, audioFileDirectory, logger)

	s := server.NewServer(router, ts)

	s.Run()
}

func getDBConnection() *sql.DB {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbServer := os.Getenv("DB_SERVER")

	dbConnStr := fmt.Sprintf("%s://%s:postgres@%s/%s?sslmode=disable", dbUser, dbPass, dbServer, dbName)

	db, err := sql.Open("postgres", dbConnStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func createTriviaService(db *sql.DB, audioFileDirectory string, logger *zap.Logger) *service.TriviaService {
	repository := postgres.NewTriviaRepository(db, logger)
	service := service.NewTriviaService(repository, audioFileDirectory, logger)
	return service
}

func newZapLogger(debug bool) (*zap.Logger, error) {
	zapLoggerConfig := zap.Config{
		Development: false,
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if debug {
		zapLoggerConfig.Development = true
		zapLoggerConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	zapLogger, err := zapLoggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to create zap logger > %w", err)
	}

	return zapLogger, nil
}
