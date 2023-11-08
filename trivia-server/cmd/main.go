package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-trivia-api/internal/db"
	"go-trivia-api/internal/server"
	"go-trivia-api/internal/service"
	"log"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-envconfig"
)

// TODO: parse all env vars in main
type config struct {
	AudioFileDirectory string `env:"AUDIO_FILE_DIRECTORY,required"`

	DatabaseName   string `env:"DB_NAME,required"`
	DatabaseUser   string `env:"DB_USER,required"`
	DatabasePass   string `env:"DB_PASS,required"`
	DatabaseServer string `env:"DB_SERVER,required"`
}

func main() {
	ctx := context.Background()

	var cfg config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	db := getDBConnection(cfg)
	runMigrations(db)

	ts := createTriviaService(cfg, db)

	s := server.NewServer(router, ts)

	s.Run()
}

func getDBConnection(cfg config) *sql.DB {
	dbConnStr := fmt.Sprintf("%s://%s:postgres@%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseServer,
		cfg.DatabaseUser)

	db, err := sql.Open("postgres", dbConnStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

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

func createTriviaService(cfg config, sqlDB *sql.DB) *service.TriviaService {
	repository := db.NewTriviaRepository(sqlDB)
	service := service.NewTriviaService(repository, cfg.AudioFileDirectory)
	return service
}
