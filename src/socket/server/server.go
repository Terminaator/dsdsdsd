package server

import (
	"log"
	"net"
)

type Server struct {
	Adr      string
	listener *net.TCPListener
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) resolve() *net.TCPAddr {
	laddr, err := net.ResolveTCPAddr("tcp", s.Adr)

	if err != nil {
		log.Fatal("Failed to resolve local address", err)
	}

	return laddr
}

func (s *Server) Listener() *net.TCPListener {
	if s.listener == nil {
		s.listen()
	}

	return s.listener
}

func (s *Server) ListenerAdr() string {
	return s.Listener().Addr().String()
}

func (s *Server) listen() {
	listener, err := net.ListenTCP("tcp", s.resolve())

	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s.listener = listener
}
