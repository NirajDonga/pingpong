package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NirajDonga/pingpong/backend/api/internal/auth"
	"github.com/NirajDonga/pingpong/backend/api/internal/config"
	"github.com/NirajDonga/pingpong/backend/api/internal/database"
	"github.com/NirajDonga/pingpong/backend/api/internal/incident"
	"github.com/NirajDonga/pingpong/backend/api/internal/middleware"
	"github.com/NirajDonga/pingpong/backend/api/internal/monitor"
	"github.com/NirajDonga/pingpong/backend/api/internal/nats"
	"github.com/NirajDonga/pingpong/backend/api/internal/result"
	"github.com/NirajDonga/pingpong/backend/api/internal/user"
	ws "github.com/NirajDonga/pingpong/backend/api/internal/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	log.Println("connecting to postgres")
	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()
	log.Println("connected to postgres")

	log.Println("tinybird configured")

	log.Println("connecting to nats")
	natsClient, err := nats.NewClient(cfg.NATSURL)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer natsClient.Close()
	log.Println("connected to nats")

	authSvc := auth.NewService(cfg.JWTSecret, 24*time.Hour)
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, authSvc)
	userHandler := user.NewHandler(userSvc, cfg.CookieSecure)
	monitorRepo := monitor.NewRepository(db)
	monitorSvc := monitor.NewService(monitorRepo)
	monitorHandler := monitor.NewHandler(monitorSvc)
	resultRepo := result.NewTinybirdRepository(cfg.TinybirdHost, cfg.TinybirdReadToken)
	resultHandler := result.NewHandler(monitorSvc, resultRepo)
	incidentRepo := incident.NewRepository(db)
	incidentHandler := incident.NewHandler(incidentRepo)
	wsManager := ws.NewManager()

	_, err = natsClient.SubscribeCheckResults(func(checkResult result.CheckResult) {
		go func() {
			wsManager.Broadcast(checkResult.MonitorID, checkResult)
		}()
	})
	if err != nil {
		log.Fatalf("check result subscription: %v", err)
	}

	router := gin.Default()
	router.Use(middleware.CORS(cfg.WebOrigin))

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "api healthy")
	})

	api := router.Group("/api")

	protected := api.Group("")
	protected.Use(middleware.Auth(authSvc))

	user.RegisterRoutes(api, protected, userHandler)
	monitor.RegisterRoutes(protected, monitorHandler)
	result.RegisterRoutes(protected, resultHandler, wsManager)
	incident.RegisterRoutes(protected, incidentHandler)

	log.Println("api service starting on :" + cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("api service failed: %v", err)
	}
}
