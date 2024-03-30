package main

import (
	"bids/routes"
	"bids/websockets"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	server := websockets.CreateServer()
	routes.AuctionRoute(app, server)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
