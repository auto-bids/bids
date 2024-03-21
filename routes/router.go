package routes

import (
	"bids/controllers"
	"github.com/gin-gonic/gin"
)

func AuctionRoute(router *gin.Engine) {
	auction := router.Group("/auction")
	{
		auction.POST("/:email", controllers.PostAuction)
		auction.GET("/:auctionid", controllers.GetAuction)
		auction.PUT("/:email/:auctionid", controllers.UpdateAuction)
	}
}
