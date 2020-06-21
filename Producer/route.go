package main

// register all route here
func (s *Server) route() {
	s.router.HandleFunc("/sendemail", s.authMiddleware(s.handler())).Methods("POST")
}
