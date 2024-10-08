package controllers

import (
	"bids/database"
	"bids/models"
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

func GetWonAuctions(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		email := c.Param("email")
		page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
		var auction []models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		timeNow := time.Now().Unix()
		stages := bson.A{
			bson.D{{"$match", bson.D{{"bidders", email}, {"end", bson.M{"$lte": timeNow}}}}},
			bson.D{
				{"$project", bson.D{
					{"_id", 1},
					{"end", 1},
					{"start", 1},
					{"owner", 1},
					{"car", 1},
					{"startFrom", 1},
					{"created", 1},
					{"minimalRise", 1},
					{"offers", bson.D{
						{"$filter", bson.D{
							{"input", "$offers"},
							{"as", "item"},
							{"cond", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$$item.offer", bson.M{"$max": "$offers.offer"}}}},
									bson.D{{"$eq", bson.A{"$$item.sender", email}}},
								}},
							}},
						}},
					}},
				}},
			},
			bson.D{{"$match", bson.D{{"offers", bson.M{"$ne": []interface{}{}}}}}},
			bson.D{{"$skip", page * 10}},
			bson.D{{"$limit", page*10 + 10}},
		}

		res, err := auctionsCollection.Aggregate(ctxDB, stages)
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
