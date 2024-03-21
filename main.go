package main

import (
	"bids/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	routes.AuctionRoute(app)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
