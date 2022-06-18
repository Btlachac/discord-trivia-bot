package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	db "github.com/Btlachac/discord-trivia-bot/postgres"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func (s *Server) handleTriviaCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			s.logger.Error("error occurred reading request body ", zap.Error(err))
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var newTrivia db.Trivia

		err = json.Unmarshal([]byte(reqBody), &newTrivia)

		if err != nil {
			s.logger.Error("error occurred parsing JSON ", zap.Error(err))
			http.Error(w, "Failed to parse JSON body", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		err = s.triviaService.AddTrivia(&ctx, &newTrivia)

		if err != nil {
			s.logger.Error("error occurred while saving trivia ", zap.Error(err))
			http.Error(w, "Failed trying to save new trivia", http.StatusInternalServerError)
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
			s.logger.Error("failed to retrieve trivia", zap.Error(err))
			http.Error(w, "Failed to retrive new trivia", http.StatusInternalServerError)
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
			s.logger.Error("failed to parse ID", zap.Error(err))
			http.Error(w, "Failed to parse id parameter", http.StatusInternalServerError)
			return
		}

		err = s.triviaService.MarkTriviaUsed(triviaId)

		if err != nil {
			s.logger.Error("failed to mark trivia as used", zap.Error(err))
			http.Error(w, "Failed to mark trivia as used", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleRoundTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roundTypes, err := s.triviaService.RoundTypesList()

		if err != nil {
			s.logger.Error("failed to retrieve round types", zap.Error(err))
			http.Error(w, "Failed to retrieve round types", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(roundTypes)
	}
}
