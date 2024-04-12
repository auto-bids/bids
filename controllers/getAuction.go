package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		auctionID := c.Param("auctionid")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		//var auction models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		id, _ := primitive.ObjectIDFromHex(auctionID)
		var auction []models.Auction
		stages := bson.A{
			bson.D{{"$match", bson.D{{"_id", id}}}},
			bson.D{{"$unwind", bson.D{{"path", "$offers"}, {"preserveNullAndEmptyArrays", true}}}},
			bson.D{{"$sort", bson.D{{"offers.offer", 1}}}},
			bson.D{{"$limit", 1}},
			bson.D{{"$group", bson.D{{"_id", "$_id"}, {"offers", bson.D{{"$push", "$offers"}}}}}},
		}
		cursor, err := auctionsCollection.Aggregate(ctxDB, stages)
		if err = cursor.All(ctxDB, &auction); err != nil {
			fmt.Println(err)
		}
		//filter := bson.D{{"_id", id}}
		//opts := options.FindOne().SetProjection(bson.D{{"_id", 1}, {"start", 1}, {"end", 1}, {"offers", bson.D{{"$elemMatch", bson.D{{"offer", bson.D{{"$max", bson.D{{"$max", "$offer"}}}}}}}}}})
		//err := auctionsCollection.FindOne(ctxDB, filter, opts).Decode(&auction)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error adding auction",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		result <- responses.Response{
			Status:  http.StatusAccepted,
			Message: "accepted",
			Data:    map[string]interface{}{"data": auction},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
