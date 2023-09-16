package main

import (
	"context"

	"github.com/Gargi2003/vyasa-cveinformation-registry/utilities"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	router := gin.Default()
	// router.GET("/api/v1/consume", handlers.ConsumeCVEInfo)
	if err := utilities.SubscribeMessage(context.Background()); err != nil {
		// Handle error
		panic(err)
	}
	router.Run(":8081")
}
