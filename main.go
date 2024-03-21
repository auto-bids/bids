package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
