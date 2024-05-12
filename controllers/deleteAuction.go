package controllers

import (
	"bids/database"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func DeleteAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		auctionid := c.Param("auctionid")
		email := c.Param("email")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		id, _ := primitive.ObjectIDFromHex(auctionid)
		filter := bson.D{{"_id", id}, {"owner", email}}
		one, err := auctionsCollection.DeleteOne(ctxDB, filter)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error removing auction",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if one.DeletedCount == 0 {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "auction not found",
				Data:    map[string]interface{}{"error": one},
			}
			return
		}
		result <- responses.Response{
			Status:  http.StatusAccepted,
			Message: "accepted",
			Data:    map[string]interface{}{"data": one},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
