package hec

import "strconv"

// Server ...
type Server struct {
	server string
	port   int
}

// NewServer ...
func NewServer(server string, port int) *Server {

	var srvr = &Server{
		server: server,
		port:   port,
	}

	return srvr
}

// GetHecPostURL ...
func (s *Server) GetHecPostURL() string {
	return "http://" + s.server + ":" + strconv.Itoa(s.port) + "/services/collector/event"
}
