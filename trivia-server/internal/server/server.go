package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"go-trivia-api/internal/db"

)

type Server struct {
	triviaService triviaService
	router        *mux.Router
}

type triviaService interface {
	//TODO
	// GetNewTrivia() (model.Trivia, error)
	AddTrivia(newTrivia db.Trivia) error
	// MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]db.RoundType, error)
}

func (s *Server) Run() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

func NewServer(router *mux.Router, triviaService triviaService) *Server {
	s := &Server{
		router:        router,
		triviaService: triviaService,
	}
	s.routes()
	return s
}
