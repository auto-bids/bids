package websockets

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				c.closeConnection()
			}
			break
		}
		c.WriteMess <- message
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
