package websockets

import (
	"bids/database"
	"bids/models"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Auction struct {
	id                  string
	currentHighestOffer models.Offer
	Clients             map[*Client]bool
	Server              *Server
	Offer               chan models.Offer
	Broadcast           chan []byte
	End                 int64
	Stop                chan bool
	AddUser             chan *Client
	RemoveUser          chan *Client
}

func CreateAuction(name string, end int64, server *Server) (*Auction, error) {
	ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	auctionCollection := database.GetCollection(database.DB, "auctions")
	id, _ := primitive.ObjectIDFromHex(name)
	var auction []models.Auction
	stages := bson.A{
		bson.D{{"$match", bson.D{{"_id", id}}}},
		bson.D{{"$unwind", bson.D{{"path", "$offers"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$sort", bson.D{{"offers.offer", -1}}}},
		bson.D{{"$limit", 1}},
		bson.D{{"$group", bson.D{{"_id", "$_id"}, {"offers", bson.D{{"$push", "$offers"}}}}}},
	}
	cursor, err := auctionCollection.Aggregate(ctxDB, stages)
	if err = cursor.All(ctxDB, &auction); err != nil {
		fmt.Println(err)
	}
	fmt.Println(auction)
	hoffer := models.Offer{}
	if len(auction[0].Offers) != 0 {
		hoffer = auction[0].Offers[0]
	}
	return &Auction{
		id:                  name,
		currentHighestOffer: hoffer,
		Clients:             make(map[*Client]bool),
		Server:              server,
		Offer:               make(chan models.Offer),
		End:                 end,
		Stop:                make(chan bool),
		AddUser:             make(chan *Client),
		RemoveUser:          make(chan *Client),
	}, nil
}
func (r *Auction) AddClient(client *Client) {
	client.WriteMess <- []byte(r.id)
	r.Clients[client] = true
}
func (r *Auction) RemoveClient(client *Client) {
	delete(r.Clients, client)
}
func (r *Auction) GetClient(client string) *Client {
	for i := range r.Clients {
		if i.UserID == client {
			return i
		}
	}
	return nil
}
func (r *Auction) endAuction() {
	message := []byte("end")
	for client := range r.Clients {
		client.WriteMess <- message
	}
}
func (r *Auction) sendOffer(offer models.Offer) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	auctionCollection := database.GetCollection(database.DB, "auctions")
	if offer.Price > r.currentHighestOffer.Price {
		r.currentHighestOffer = offer
		id, _ := primitive.ObjectIDFromHex(r.id)
		filter := bson.D{{"_id", id}}
		update := bson.M{"$push": bson.M{"offers": offer}}
		_, err := auctionCollection.UpdateOne(ctx, filter, update)
		if err == nil {
			data, _ := json.Marshal(offer)
			for client := range r.Clients {
				client.WriteMess <- data
			}
		} else {

			r.GetClient(offer.Sender).WriteMess <- []byte("error")
		}

	}

}
func (r *Auction) RunAuction() {

	for {
		if time.Now().Unix() == r.End {
			r.endAuction()
			return
		}
		select {
		case offer := <-r.Offer:
			r.sendOffer(offer)
		case user := <-r.AddUser:
			r.AddClient(user)
		case <-r.Stop:
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}
