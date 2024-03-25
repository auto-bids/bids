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

func UpdateAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		auctionid := c.Param("auctionid")
		email := c.Param("email")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		var res models.UpdateAuction
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := c.ShouldBindJSON(&res)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		err = validate.Struct(res)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "validation failed",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		filter := bson.D{{"_id", auctionid}, {"owner", email}}
		auctionsCollection.UpdateOne(ctxDB, filter, res)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error updating auction",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		result <- responses.Response{
			Status:  http.StatusAccepted,
			Message: "accepted",
			Data:    map[string]interface{}{"data": res},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
