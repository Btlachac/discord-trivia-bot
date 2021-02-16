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

// TODO: mark trivia as used

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbServer := os.Getenv("DB_SERVER")

	dbConnStr := fmt.Sprintf("%s://%s:postgres@%s/%s?sslmode=disable", dbUser, dbPass, dbServer, dbName)

	models.DB, err = sql.Open("postgres", dbConnStr)

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/trivia", createTrivia).Methods("POST")
	router.HandleFunc("/trivia", getNewTrivia).Methods("GET")
	router.HandleFunc("/trivia/{id:[0-9]+}/mark-used", markTriviaUsed).Methods("PUT")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

func createTrivia(w http.ResponseWriter, r *http.Request) {

	// TODO: handle incoming files
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	var newTrivia models.Trivia

	// fmt.Println(string(reqBody))

	json.Unmarshal([]byte(reqBody), &newTrivia)

	models.AddTrivia(newTrivia)

	// fmt.Printf("%+v\n", newTrivia)
	// events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTrivia)
}

func getNewTrivia(w http.ResponseWriter, r *http.Request) {

	trivia := models.GetNewTrivia()

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(trivia)

}

func markTriviaUsed(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	triviaId, _ := strconv.ParseInt(params["id"], 10, 64)

	models.MarkTriviaUsed(triviaId)
}
