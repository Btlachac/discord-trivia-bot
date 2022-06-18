package server

import (
	"log"
	"net/http"

	"github.com/Btlachac/discord-trivia-bot/service"
	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

type Server struct {
	triviaService service.Service
	router        *mux.Router
}

func (s *Server) Run() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

func NewServer(router *mux.Router, triviaService service.Service) *Server {
	s := &Server{
		router:        router,
		triviaService: triviaService,
	}
	s.routes()
	return s
}
