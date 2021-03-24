package server

import (
	"encoding/json"
	"go-trivia-api/model"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) handleTriviaCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var newTrivia model.Trivia

		json.Unmarshal([]byte(reqBody), &newTrivia)

		err = s.triviaService.AddTrivia(newTrivia)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newTrivia)
	}
}

func (s *Server) handleTriviaGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		trivia, err := s.triviaService.GetNewTrivia()

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(trivia)
	}
}

func (s *Server) handleTriviaMarkUsed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		triviaId, err := strconv.ParseInt(params["id"], 10, 64)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.triviaService.MarkTriviaUsed(triviaId)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
