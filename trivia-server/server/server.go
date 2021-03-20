package server

import (
	"database/sql"
	"fmt"
	"go-trivia-api/repository"
	"go-trivia-api/service"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Server struct {
	triviaService *service.TriviaService
	router        *mux.Router
}


func Run() {

	router := mux.NewRouter().StrictSlash(true)

	s := newServer(router)

	db := getDBConnection()

	runMigrations(db)

	s.triviaService = createTriviaService(db)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

func createTriviaService(db *sql.DB) *service.TriviaService {
	repository := repository.NewTriviaRepository(db)
	service := service.NewTriviaService(repository)
	return service
}

func newServer(router *mux.Router) *Server {
	s := &Server{}
	s.router = router
	s.routes()
	return s
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
	driver, err := postgres.WithInstance(db, &postgres.Config{})

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
