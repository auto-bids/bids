package controllers

import (
	"bids/database"
	"bids/models"
	"bids/queries"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

func GetAllAuctions(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		page, _ := strconv.ParseInt(c.Param("page"), 10, 64)
		order, _ := strconv.ParseInt(c.Query("order"), 10, 8)
		by := c.Query("sortby")
		sort := models.Sort{Order: order, By: by}
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(sort); err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error validation sort query",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		var car models.CarSearch
		err := c.ShouldBindJSON(&car)
		if err != nil {
			return
		}
		filter := queries.GetOfferQuery(car)

		var auction []models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		opts := options.Find().SetLimit(page * 10)
		opts.SetSort(bson.D{{"car." + sort.By, sort.Order}})
		res, err := auctionsCollection.Find(ctxDB, filter, opts)
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
				Message: "Error decoding auctions",
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
