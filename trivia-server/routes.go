package main

func (s *server) routes() {
	s.router.HandleFunc("/trivia", s.handleTriviaCreate()).Methods("POST")
	s.router.HandleFunc("/trivia", s.handleTriviaGet()).Methods("GET")
	s.router.HandleFunc("/trivia/{id:[0-9]+}/mark-used", s.handleTriviaMarkUsed()).Methods("PUT")
}
