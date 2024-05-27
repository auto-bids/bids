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
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

func GetAllAuctions(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		page, _ := strconv.ParseInt(c.Param("page"), 10, 64)
		order := c.Query("order")
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
		err := c.ShouldBind(&car)
		if err != nil {
			return
		}
		filter := queries.GetOfferQuery(car)
		var auction []models.Auction
		auctionsCollection := database.GetCollection(database.DB, "auctions")
		matchStage := bson.D{{"$match", filter}}
		sortStage := bson.D{{"$sort", bson.D{{"car.title", 1}}}}
		projectStage := bson.D{
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
						{"cond", bson.D{{"$eq", bson.A{"$$item.offer", bson.M{"$max": "$offers.offer"}}}}}},
					},
				}},
			}},
		}
		skipStage := bson.D{{"$skip", page * 10}}
		limitStage := bson.D{{"$limit", page*10 + 10}}
		if sort.By != "" && sort.Order != "" {
			var orderi int8
			switch order {
			case "desc":
				orderi = -1
			case "asc":
				orderi = 1
			default:
				orderi = -1
			}
			sortStage = bson.D{{"$sort", bson.D{{"car." + sort.By, orderi}}}}
		}
		stages := mongo.Pipeline{matchStage, sortStage, projectStage, skipStage, limitStage}
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

