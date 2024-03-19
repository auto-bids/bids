package controlers

import (
	"bids/responses"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func postAuction(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		email := ctx.Param("email")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
	}(ctx.Copy())
}
