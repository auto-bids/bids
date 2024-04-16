package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

func GetOffers(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {

		page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		roomCollection := database.GetCollection(database.DB, "auctions")
		id, err := primitive.ObjectIDFromHex(c.Param("auctionid"))
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		filterUser := bson.D{{"_id", id}}
		err = roomCollection.FindOne(ctxDB, filterUser).Err()
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		pipeline := bson.A{
			bson.D{{"$match", bson.D{{"_id", id}}}},
			bson.D{{"$unwind", "$offers"}},
			bson.D{{"$sort", bson.D{{"offers.time", -1}}}},
			bson.D{{"$skip", page * 10}},
			bson.D{{"$limit", 10}},
		}
		cursor, err := roomCollection.Aggregate(ctxDB, pipeline)
		var results []models.OfferUnwind
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "auction not found",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if err = cursor.All(ctxDB, &results); err != nil {
			panic(err)
		}

		result <- responses.Response{
			Status:  http.StatusFound,
			Message: "auction found",
			Data:    map[string]interface{}{"data": results},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
