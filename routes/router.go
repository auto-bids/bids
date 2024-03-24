package routes

import (
	"bids/controllers"
	"github.com/gin-gonic/gin"
)

func AuctionRoute(router *gin.Engine) {
	auction := router.Group("/auction")
	{
		auction.POST("/addBid/:email", controllers.PostAuction)
		auction.GET("/getBid/:auctionid", controllers.GetAuction)
		auction.DELETE("/removeBid/:email/:auctionid", controllers.DeleteAuction)
		auction.PUT("/editBid/:email/:auctionid", controllers.UpdateAuction)
	}
}
