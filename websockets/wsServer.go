package websockets

import (
	"sync"
)

type Server struct {
	Mutex   sync.Mutex
	Clients map[*Client]bool
}

func CreateServer() *Server {
	return &Server{
		Clients: make(map[*Client]bool),
	}
}

func (s *Server) AddClient(client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[client] = true
}
