package main

import (
	"database/sql"
	"fmt"
	"go-trivia-api/postgres"
	"go-trivia-api/server"
	"go-trivia-api/service"
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

func createTriviaService(db *sql.DB) *service.TriviaService {
	repository := postgres.NewTriviaRepository(db)
	service := service.NewTriviaService(repository)
	return service
}
