// Package main is the entry point for the IoT backend API.
//
// @title                      IoT Backend API
// @version                    1.0
// @description                CRUD for User & Device Access Control (Echo + PostgreSQL, no ORM).
// @BasePath                   /api/v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/vinzryyy/iot-backend/app"
	"github.com/vinzryyy/iot-backend/database"
	_ "github.com/vinzryyy/iot-backend/docs"
	"github.com/vinzryyy/iot-backend/repo"
	"github.com/vinzryyy/iot-backend/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	cfg := app.LoadConfig()

	pool, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("database connect: %v", err)
	}
	defer pool.Close()

	if err := database.Migrate(pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := database.SeedUsers(pool); err != nil {
		log.Fatalf("seed users: %v", err)
	}

	// Repositories
	userRepo := repo.NewUserRepository(pool)
	deviceRepo := repo.NewDeviceRepository(pool)
	locationRepo := repo.NewLocationRepository(pool)

	// Services
	jwtSvc := service.NewJWTService(cfg.JWTSecret, cfg.JWTExpHours)
	authSvc := service.NewAuthService(userRepo, locationRepo, jwtSvc)
	deviceSvc := service.NewDeviceService(deviceRepo, userRepo, locationRepo)

	// Handlers
	authHandler := app.NewAuthHandler(authSvc)
	deviceHandler := app.NewDeviceHandler(deviceSvc)

	e := app.NewRouter(app.Deps{
		AuthHandler:   authHandler,
		DeviceHandler: deviceHandler,
		JWT:           jwtSvc,
	})

	// Run server with graceful shutdown
	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		log.Printf("server listening on %s (env=%s)", addr, cfg.AppEnv)
		if err := e.Start(addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
