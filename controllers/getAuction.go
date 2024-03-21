package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func GetAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		auctionID := ctx.Param("auctionid")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		var auction models.Auction
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(auction)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "validation failed",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		filter := bson.D{{"_id", auctionID}}
		auctionsCollection.FindOne(ctxDB, filter).Decode(&auction)
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
