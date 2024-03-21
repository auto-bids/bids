package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

func PostAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		email := ctx.Param("email")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		var res models.Auction
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
		res.Owner = email
		res.Created = time.Now()
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
		one, err := auctionsCollection.InsertOne(ctxDB, res)
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
			Data:    map[string]interface{}{"data": one},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
