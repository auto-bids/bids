package websockets

import (
	"bids/database"
	"bids/models"
	"bids/responses"
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
	minimalRaise        int64
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
		bson.D{
			{"$project", bson.D{
				{"_id", "$_id"},
				{"minimalRaise", "$minimalRaise"},
				{"offers", bson.D{
					{"$filter", bson.D{
						{"input", "$offers"},
						{"as", "item"},
						{"cond", bson.D{
							{"$eq", bson.A{"$$item.offer", bson.M{"$max": "$offers.offer"}}},
						},
						}},
					}},
				}},
			}},
	}
	cursor, err := auctionCollection.Aggregate(ctxDB, stages)
	if err = cursor.All(ctxDB, &auction); err != nil {
		fmt.Println(err)
	}
	hoffer := models.Offer{}
	if len(auction[0].Offers) != 0 {
		hoffer = auction[0].Offers[0]
	}
	minimalRaise := auction[0].MinimalRaise
	return &Auction{
		id:                  name,
		currentHighestOffer: hoffer,
		minimalRaise:        minimalRaise,
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
	fmt.Println()
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
	message, err := json.Marshal(r.currentHighestOffer)
	if err != nil {
		for client := range r.Clients {
			client.WriteMess <- []byte("end")
		}
	}
	for client := range r.Clients {
		client.WriteMess <- message
	}
}
func (r *Auction) sendOffer(offer models.Offer) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	auctionCollection := database.GetCollection(database.DB, "auctions")
	if offer.Price > r.currentHighestOffer.Price+r.minimalRaise {
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

	} else {
		mess, _ := json.Marshal(responses.ResponseWs{Message: "offer too low"})
		r.GetClient(offer.Sender).WriteMess <- mess
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
		case user := <-r.RemoveUser:
			r.RemoveClient(user)
			if len(r.Clients) == 0 {
				r.Server.RemoveAuction(r.id)
				return
			}
		case <-r.Stop:
			return
		default:
		}
		time.Sleep(time.Microsecond * 500)
	}

}
