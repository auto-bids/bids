package routes

import (
	"bids/controllers"
	"github.com/gin-gonic/gin"
)

func auctionRoute(router *gin.Engine) {
	auction := router.Group("/auction")
	{
		auction.POST("/:email", controllers.PostAuction)
	}
}
