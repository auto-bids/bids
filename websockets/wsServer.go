package websockets

import (
	"sync"
)

type Server struct {
	Mutex   sync.Mutex
	Clients map[*Client]bool
	Rooms   map[string]*Auction
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
func (s *Server) AddAuction(id string) (*Auction, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	auct, err := CreateAuction(id, s)
	if err != nil {
		return nil, err
	}

	go auct.RunAuction()
	s.Rooms[id] = auct
	return auct, nil
}
