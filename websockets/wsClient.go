package websockets

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Client struct {
	Socket    *websocket.Conn
	WriteMess chan []byte
	Server    *Server
	UserID    string
	Close     chan string
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func NewClient(socket *websocket.Conn, ctx *gin.Context) *Client {

	return &Client{
		Socket:    socket,
		Close:     make(chan string),
		WriteMess: make(chan []byte),
		UserID:    ctx.Param("email"),
	}
}
func (c *Client) JoinAuction(dest string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	auctionCollection := database.GetCollection(database.DB, "auctions")
	id, _ := primitive.ObjectIDFromHex(dest)
	filter := bson.D{{"_id", id}}
	var auction models.GetAuctionForRoom
	err := auctionCollection.FindOne(ctx, filter).Decode(&auction)
	if err != nil {
		wsErr := responses.ResponseWs{
			Message: "auction not found",
			Data:    map[string]interface{}{"error": err},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	if auction.End <= time.Now().Unix() {
		wsErr := responses.ResponseWs{
			Message: "auction has ended",
			Data:    map[string]interface{}{"error": err},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	update := bson.M{"$push": bson.M{"users": c.UserID}}
	updateRes, err := auctionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		wsErr := responses.ResponseWs{
			Message: "user not added",
			Data:    map[string]interface{}{"error": err},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	wsRes := responses.ResponseWs{
		Message: "user added",
		Data:    map[string]interface{}{"data": updateRes},
	}
	res, _ := json.Marshal(wsRes)
	c.WriteMess <- res
}
func (c *Client) makeBid(offer models.Offer) {
	offer.Sender = c.UserID
	offer.Time = time.Now().UnixNano()

}
func (c *Client) closeConnection() {
	delete(c.Server.Clients, c)
}
func (c *Client) ReadPump() {
	defer func() {
		c.Socket.Close()
	}()
	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error { c.Socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, byteMessage, err := c.Socket.ReadMessage()
		mess := &models.Message{}
		json.Unmarshal(byteMessage, mess)
		mess.Sender = c.UserID
		switch mess.Options {
		case "join":
			c.JoinAuction(mess.Destination)
		case "bid":
			//c.makeBid(mess) TODO
		}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				c.closeConnection()
			}
			break
		}

		time.Sleep(time.Millisecond)
	}
}
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.WriteMess:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.WriteMess)
			for i := 0; i < n; i++ {
				w.Write(<-c.WriteMess)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
		time.Sleep(time.Millisecond)
	}
}
