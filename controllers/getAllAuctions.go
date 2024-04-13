package controllers

import (
	"bids/database"
	"bids/models"
	"bids/queries"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetAllAuctions(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		//page := c.Param("page")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var car models.CarSearch
		err := c.ShouldBindJSON(&car)
		if err != nil {
			return
		}
		filter := queries.GetOfferQuery(car)
		defer close(result)
		defer cancel()
		var auction []models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")

		res, err := auctionsCollection.Find(ctxDB, filter)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error adding auction",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if err := res.All(ctx, &auction); err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error decoding offers",
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
