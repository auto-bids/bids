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
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Client struct {
	Socket    *websocket.Conn
	WriteMess chan []byte
	Server    *Server
	Auctions  map[string]*Auction
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
		Auctions:  make(map[string]*Auction),
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
	err := auctionCollection.FindOne(ctx, filter, options.FindOne().SetProjection(bson.D{{"end", 1}, {"start", 1}})).Decode(&auction)
	if err != nil {
		wsErr := responses.ResponseWs{
			Message: "auction not found",
			Data:    map[string]interface{}{"error": err},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	if auction.End < time.Now().Unix() {
		wsErr := responses.ResponseWs{
			Message: "ended",
			Data:    map[string]interface{}{"error": "ended"},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	if auction.Start > time.Now().Unix() {
		wsErr := responses.ResponseWs{
			Message: "auction has not started yet",
			Data:    map[string]interface{}{"error": "notStarted"},
		}
		res, _ := json.Marshal(wsErr)
		c.WriteMess <- res
		return
	}
	update := bson.M{"$addToSet": bson.M{"bidders": c.UserID}}
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
	auctionServer := c.Server.GetAuction(dest)
	if auctionServer == nil {
		c.Server.AddAuction(dest, auction.End)
		auctionServer = c.Server.GetAuction(dest)
	}
	c.Auctions[auctionServer.id] = auctionServer
	auctionServer.AddUser <- c
	wsRes := responses.ResponseWs{
		Message: "user added",
		Data:    map[string]interface{}{"data": updateRes},
	}
	res, _ := json.Marshal(&wsRes)
	c.WriteMess <- res
}
func (c *Client) LeaveAuction(dest string) {
	auct := c.Auctions[dest]
	mess, _ := json.Marshal(responses.ResponseWs{Message: "disconnected"})
	if auct != nil {
		auct.RemoveUser <- c
	}
	c.WriteMess <- mess
}
func (c *Client) makeBid(mess *models.Message) {
	offer := mess.Offer
	offer.Sender = c.UserID
	offer.Time = time.Now().UnixNano()
	if c.Auctions[mess.Destination] != nil {
		c.Auctions[mess.Destination].Offer <- offer
	}
}
func (c *Client) closeConnection() {
	for _, v := range c.Auctions {
		v.RemoveUser <- c
	}
	c.Server.RemoveClient(c.UserID)
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
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				c.closeConnection()
			}
			break
		}
		mess := &models.Message{}
		json.Unmarshal(byteMessage, mess)
		mess.Sender = c.UserID
		switch mess.Options {
		case "join":
			c.JoinAuction(mess.Destination)
		case "leave":
			c.LeaveAuction(mess.Destination)
		case "bid":
			c.makeBid(mess)
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
