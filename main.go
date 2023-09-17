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

	if err := utilities.SubscribeMessage(context.Background()); err != nil {
		panic(err)
	}
	router.Run(":8081")
}
