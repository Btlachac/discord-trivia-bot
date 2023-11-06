package server

import (
	"go-trivia-api/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

type Server struct {
	triviaService triviaService
	router        *mux.Router
}

type triviaService interface {
	//TODO
	// GetNewTrivia() (model.Trivia, error)
	AddTrivia(newTrivia model.Trivia) error
	// MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]model.RoundType, error)
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
