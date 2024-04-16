package routes

import (
	"bids/controllers"
	"bids/websockets"
	"github.com/gin-gonic/gin"
)

func AuctionRoute(router *gin.Engine, Server *websockets.Server) {
	auction := router.Group("/auction")
	{
		auction.POST("/add/:email", controllers.PostAuction)
		auction.GET("/get/:auctionid", controllers.GetAuction)
		auction.GET("/my/:email", controllers.GetUserAuctions)
		auction.GET("/auction/:auctionid/:page", controllers.GetOffers)
		auction.GET("/won/:email", controllers.GetWonAuctions)
		auction.GET("/joined/:email", controllers.GetJoinedAuctions)
		auction.GET("/search/:page", controllers.GetAllAuctions)
		auction.DELETE("/remove/:email/:auctionid", controllers.DeleteAuction)
		auction.PUT("/edit/:email/:auctionid", controllers.UpdateAuction)
		auction.GET("/ws/:email", func(ctx *gin.Context) { websockets.ManageWs(Server, ctx) })
	}
}
