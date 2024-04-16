package websockets

import (
	"sync"
)

type Server struct {
	Mutex    sync.Mutex
	Clients  map[string]*Client
	Auctions map[string]*Auction
}

func CreateServer() *Server {
	return &Server{
		Clients:  make(map[string]*Client),
		Auctions: make(map[string]*Auction),
	}
}
func (s *Server) RemoveAuction(id string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Auctions, id)
}
func (s *Server) AddClient(client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[client.UserID] = client
}
func (s *Server) RemoveClient(client string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Clients, client)
}
func (s *Server) GetAuction(id string) *Auction {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if s.Auctions[id] == nil {
		return nil
	}
	return s.Auctions[id]
}
func (s *Server) AddAuction(id string, end int64) (*Auction, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	auct, err := CreateAuction(id, end, s)
	if err != nil {
		return nil, err
	}
	go auct.RunAuction()
	s.Auctions[id] = auct
	return auct, nil
}
