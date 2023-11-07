package server

import (
	"encoding/json"
	"go-trivia-api/internal/db"
	"io"
	"log"
	"net/http"
)

func (s *Server) handleTriviaCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)

		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var newTrivia db.Trivia

		err = json.Unmarshal([]byte(reqBody), &newTrivia)

		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to parse JSON body", http.StatusInternalServerError)
			return
		}

		err = s.triviaService.AddTrivia(r.Context(), newTrivia)

		if err != nil {
			log.Print(err)
			http.Error(w, "Failed trying to save new trivia", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(newTrivia); err != nil {
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
			log.Print(err)
			http.Error(w, "Failed to retrieve round types", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(roundTypes); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
