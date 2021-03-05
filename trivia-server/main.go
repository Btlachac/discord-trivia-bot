package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-trivia-api/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type server struct {
	db     *sql.DB
	router *mux.Router
	model  *models.Model
}

func newServer(db *sql.DB, router *mux.Router, model *models.Model) *server {
	s := &server{}
	s.db = db
	s.router = router
	s.model = model
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := getDBConnection()

	model := models.NewModel(db)

	router := mux.NewRouter().StrictSlash(true)

	s := newServer(db, router, model)

	s.routes()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	log.Fatal(http.ListenAndServe(":8080", handler))

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

func (s *server) handleTriviaCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatal(err)
		}

		var newTrivia models.Trivia

		json.Unmarshal([]byte(reqBody), &newTrivia)

		s.model.AddTrivia(newTrivia)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newTrivia)
	}
}

func (s *server) handleTriviaGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		trivia := s.model.GetNewTrivia()

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(trivia)
	}
}

func (s *server) handleTriviaMarkUsed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		triviaId, _ := strconv.ParseInt(params["id"], 10, 64)

		s.model.MarkTriviaUsed(triviaId)
	}
}
