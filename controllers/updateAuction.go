package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		id, _ := primitive.ObjectIDFromHex(auctionid)
		var res models.UpdateAuction
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := c.ShouldBindJSON(&res)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request",
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
		var check models.GetAuctionShort
		filter := bson.D{{"_id", id}, {"owner", email}}
		auctionsCollection.FindOne(ctxDB, filter).Decode(&check)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "offer not found",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if time.Now().Unix() >= check.Start {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "can't update an auction",
				Data:    map[string]interface{}{},
			}
			return
		}

		update := bson.D{{"$set", bson.D{
			{"minimalRaise", res.MinimalRaise},
			{"car.title", res.Title},
			{"car.description", res.Description},
			{"car.photos", res.Photos},
			{"car.year", res.Year},
			{"car.mileage", res.Mileage},
			{"car.telephone_number", res.TelephoneNumber}}}}
		one, err := auctionsCollection.UpdateOne(ctxDB, filter, update)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error updating auction",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
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
			Data:    map[string]interface{}{"data": one},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
