package server

func (s *Server) routes() {
	s.router.HandleFunc("/trivia", s.handleTriviaCreate()).Methods("POST")
	//TODO 
	// s.router.HandleFunc("/trivia", s.handleTriviaGet()).Methods("GET")
	// s.router.HandleFunc("/trivia/{id:[0-9]+}/mark-used", s.handleTriviaMarkUsed()).Methods("PUT")
	s.router.HandleFunc("/trivia/roundTypes", s.handleRoundTypes()).Methods("GET")
}
