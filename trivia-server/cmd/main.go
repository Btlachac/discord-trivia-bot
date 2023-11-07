package main

import (
	"database/sql"
	"fmt"
	"go-trivia-api/internal/db"
	"go-trivia-api/internal/server"
	"go-trivia-api/internal/service"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	db := getDBConnection()
	runMigrations(db)

	ts := createTriviaService(db)

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

func createTriviaService(sqlDB *sql.DB) *service.TriviaService {
	repository := db.NewTriviaRepository(sqlDB)
	service := service.NewTriviaService(repository)
	return service
}
