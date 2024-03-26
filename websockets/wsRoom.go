package websockets

import (
	"bids/database"
	"bids/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Room struct {
	id         string
	Clients    map[*Client]bool
	Server     *Server
	Broadcast  chan *models.Message
	End        int64
	Stop       chan bool
	AddUser    chan *Client
	RemoveUser chan *Client
}

func CreateRoom(name string, end int64, server *Server) *Room {
	return &Room{
		id:         name,
		Server:     server,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *models.Message),
		End:        end,
		AddUser:    make(chan *Client),
		RemoveUser: make(chan *Client),
		Stop:       make(chan bool),
	}
}
func (r *Room) AddClient(client *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roomCollection := database.GetCollection(database.DB, "rooms")
	id, _ := primitive.ObjectIDFromHex(r.id)
	filter := bson.D{{"_id", id}, {"users", client.UserID}}
	var room models.RoomDB
	err := roomCollection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		client.WriteMess <- []byte("unauthorized")
		return
	}
	client.WriteMess <- []byte(r.id)
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client) {
	delete(r.Clients, client)
}
func (r *Room) GetClient(client string) *Client {
	for i := range r.Clients {
		if i.UserID == client {
			return i
		}
	}
	return nil
}
func (r *Room) endAuction() {
	message := []byte("end")
	for client := range r.Clients {
		client.WriteMess <- message
	}
}
func (r *Room) sendMessage(message *models.Message) {
	//TODO
}
func (r *Room) RunRoom() {

	for {
		if time.Now().Unix() == r.End {
			r.endAuction()
			return
		}
		select {
		case message := <-r.Broadcast:
			r.sendMessage(message)
		case user := <-r.AddUser:
			r.AddClient(user)
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
			if len(r.Clients) == 0 {
				return
			}
		case <-r.Stop:
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}
