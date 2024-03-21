package websockets

import (
	"bids/responses"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ManageWs(server *Server, ctx *gin.Context) {

	ws, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	client := NewClient(ws, ctx)
	server.AddClient(client)
	if err != nil {
		connectionError := responses.Response{
			Status:  http.StatusBadRequest,
			Message: "websocket connectin failed",
			Data:    map[string]interface{}{"err": err.Error()},
		}
		ctx.JSON(connectionError.Status, connectionError)
		return
	}

}
