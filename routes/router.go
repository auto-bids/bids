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
		auction.GET("/my/:email/:page", controllers.GetUserAuctions)
		auction.GET("/offers/:auctionid/:page", controllers.GetOffers)
		auction.GET("/won/:email/:page", controllers.GetWonAuctions)
		auction.GET("/joined/:email/:page", controllers.GetJoinedAuctions)
		auction.GET("/search/:page", controllers.GetAllAuctions)
		auction.DELETE("/remove/:email/:auctionid", controllers.DeleteAuction)
		auction.PUT("/edit/:email/:auctionid", controllers.UpdateAuction)
		auction.GET("/ws/:email", func(ctx *gin.Context) { websockets.ManageWs(Server, ctx) })
	}
}
