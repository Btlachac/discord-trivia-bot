package server

func (s *Server) routes() {
	s.router.HandleFunc("/trivia", s.handleTriviaCreate()).Methods("POST")
	s.router.HandleFunc("/trivia", s.handleTriviaGet()).Methods("GET")
	s.router.HandleFunc("/trivia/list", s.handleTriviaList()).Methods("GET")
	s.router.HandleFunc("/trivia/{id:[0-9]+}/mark-used", s.handleTriviaMarkUsed()).Methods("PUT")
}
