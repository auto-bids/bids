package websockets

import (
	"fmt"
	"sync"
)

type Server struct {
	Mutex    sync.Mutex
	Clients  map[*Client]bool
	Auctions map[string]*Auction
}

func CreateServer() *Server {
	return &Server{
		Clients:  make(map[*Client]bool),
		Auctions: make(map[string]*Auction),
	}
}

func (s *Server) AddClient(client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[client] = true
	fmt.Println(s.Clients[client])
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
	s.Auctions[id] = auct
	go auct.RunAuction()
	return auct, nil
}
