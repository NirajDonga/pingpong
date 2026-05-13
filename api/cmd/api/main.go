package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/NirajDonga/pingpong/api/internal/config"
	"github.com/NirajDonga/pingpong/api/internal/database"
	"github.com/NirajDonga/pingpong/api/internal/incident"
	"github.com/NirajDonga/pingpong/api/internal/middleware"
	"github.com/NirajDonga/pingpong/api/internal/monitor"
	"github.com/NirajDonga/pingpong/api/internal/nats"
	"github.com/NirajDonga/pingpong/api/internal/result"
	"github.com/NirajDonga/pingpong/api/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()

	if err := database.PingClickHouse(context.Background(), cfg.ClickHouseURL); err != nil {
		log.Fatalf("clickhouse connect: %v", err)
	}

	natsClient, err := nats.NewClient(cfg.NATSURL)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer natsClient.Close()

	authSvc := auth.NewService(cfg.JWTSecret, 24*time.Hour)
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, authSvc)
	userHandler := user.NewHandler(userSvc)
	monitorRepo := monitor.NewRepository(db)
	monitorSvc := monitor.NewService(monitorRepo)
	monitorHandler := monitor.NewHandler(monitorSvc)
	resultRepo := result.NewClickHouseRepository(cfg.ClickHouseURL)
	resultSvc := result.NewService(resultRepo, monitorRepo)
	resultHandler := result.NewHandler(monitorSvc, resultSvc)
	incidentRepo := incident.NewRepository(db)
	incidentSvc := incident.NewService(incidentRepo)
	incidentHandler := incident.NewHandler(incidentSvc)

	_, err = natsClient.SubscribeCheckResults(func(checkResult result.CheckResult) {
		go func() {
			if err := resultSvc.Process(context.Background(), checkResult); err != nil {
				log.Printf("failed to process check result for monitor %s: %v", checkResult.MonitorID, err)
				return
			}

			log.Printf("stored check result for monitor %s success=%t status=%d", checkResult.MonitorID, checkResult.Success, checkResult.StatusCode)
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
	result.RegisterRoutes(protected, resultHandler)
	incident.RegisterRoutes(protected, incidentHandler)

	log.Println("api service starting on :" + cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("api service failed: %v", err)
	}
}
