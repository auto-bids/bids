package routes

import (
	"bids/controllers"
	"bids/websockets"
	"github.com/gin-gonic/gin"
)

func AuctionRoute(router *gin.Engine, Server *websockets.Server) {
	auction := router.Group("/auction")
	{
		auction.POST("/addBid/:email", controllers.PostAuction)
		auction.GET("/getBid/:auctionid", controllers.GetAuction)
		auction.GET("/auctions/:email", controllers.GetUserAuctions)
		auction.GET("/auctionsWon/:email", controllers.GetWonAuctions)
		auction.GET("/auctionsJoined/:email", controllers.GetJoinedAuctions)
		auction.GET("/search/:page", controllers.GetAllAuctions)
		auction.DELETE("/removeBid/:email/:auctionid", controllers.DeleteAuction)
		auction.PUT("/editBid/:email/:auctionid", controllers.UpdateAuction)
		auction.GET("/ws/:email", func(ctx *gin.Context) { websockets.ManageWs(Server, ctx) })
	}
}
