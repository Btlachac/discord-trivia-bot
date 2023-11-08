package server

import (
	"encoding/json"
	"go-trivia-api/internal/db"
	"io"
	"log/slog"
	"net/http"
)

func (s *Server) handleTriviaCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)

		if err != nil {
			slog.Error("error reading request body", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var newTrivia db.Trivia

		err = json.Unmarshal([]byte(reqBody), &newTrivia)

		if err != nil {
			slog.Error("failed to unmarshal request", err)
			http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
			return
		}

		err = s.triviaService.AddTrivia(r.Context(), newTrivia)

		if err != nil {
			slog.Error("failed saving new trivia", err)
			http.Error(w, "Failed trying to save new trivia", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(newTrivia); err != nil {
			slog.Error("failed encoding response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

//TODO
// func (s *Server) handleTriviaGet() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		trivia, err := s.triviaService.GetNewTrivia()

// 		if err != nil {
// 			log.Print(err)
// 			http.Error(w, "Failed to retrive new trivia", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)

// 		json.NewEncoder(w).Encode(trivia)
// 	}
// }

// func (s *Server) handleTriviaMarkUsed() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		params := mux.Vars(r)

// 		triviaId, err := strconv.ParseInt(params["id"], 10, 64)

// 		if err != nil {
// 			log.Print(err)
// 			http.Error(w, "Failed to parse id parameter", http.StatusInternalServerError)
// 			return
// 		}

// 		err = s.triviaService.MarkTriviaUsed(triviaId)

// 		if err != nil {
// 			log.Print(err)
// 			http.Error(w, "Failed to mark trivia as used", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}
// }

func (s *Server) handleRoundTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roundTypes, err := s.triviaService.RoundTypesList(r.Context())

		if err != nil {
			slog.Error("error retrieving round types", err)
			http.Error(w, "Failed to retrieve round types", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(roundTypes); err != nil {
			slog.Error("failed to encode response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
