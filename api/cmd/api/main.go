package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "api healthy")
	})

	log.Println("api service starting on :3001")

	if err := router.Run(":3001"); err != nil {
		log.Fatalf("api service failed: %v", err)
	}
}
