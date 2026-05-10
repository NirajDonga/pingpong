package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/NirajDonga/pingpong/api/internal/config"
	"github.com/NirajDonga/pingpong/api/internal/database"
	"github.com/NirajDonga/pingpong/api/internal/middleware"
	"github.com/NirajDonga/pingpong/api/internal/monitor"
	"github.com/NirajDonga/pingpong/api/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()

	authSvc := auth.NewService(cfg.JWTSecret, 24*time.Hour)
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, authSvc)
	userHandler := user.NewHandler(userSvc)
	monitorRepo := monitor.NewRepository(db)
	monitorSvc := monitor.NewService(monitorRepo)
	monitorHandler := monitor.NewHandler(monitorSvc)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "api healthy")
	})

	api := router.Group("/api")

	protected := api.Group("")
	protected.Use(middleware.Auth(authSvc))

	user.RegisterRoutes(api, protected, userHandler)
	monitor.RegisterRoutes(protected, monitorHandler)

	log.Println("api service starting on :" + cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("api service failed: %v", err)
	}
}
