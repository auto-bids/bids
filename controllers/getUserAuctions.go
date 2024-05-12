package controllers

import (
	"bids/database"
	"bids/models"
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

func GetUserAuctions(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		email := c.Param("email")
		page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
		status := models.Status{Status: c.Query("status")}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(status); err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error validation sort query",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		filter := bson.M{"owner": email}
		if status.Status != "" {
			switch status.Status {
			case "ended":
				filter["end"] = bson.M{"$lte": time.Now().Unix()}
			case "ongoing":
				filter["end"] = bson.M{"$gte": time.Now().Unix()}
				filter["start"] = bson.M{"$lte": time.Now().Unix()}
			case "notstarted":
				filter["start"] = bson.M{"$gte": time.Now().Unix()}
			}
		}
		var auction []models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		opts := options.Find().SetSkip(page * 10).SetLimit(page*10 + 10)
		res, err := auctionsCollection.Find(ctxDB, filter, opts)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error searching auction",
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
