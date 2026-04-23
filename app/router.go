package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/vinzryyy/iot-backend/service"
)

type Deps struct {
	AuthHandler   *AuthHandler
	DeviceHandler *DeviceHandler
	JWT           *service.JWTService
}

func NewRouter(d Deps) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Validator = NewValidator()
	e.HTTPErrorHandler = ErrorHandler

	// Global middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","id":"${id}","method":"${method}",` +
			`"uri":"${uri}","status":${status},"latency":"${latency_human}"}` + "\n",
	}))
	e.Use(PrometheusMiddleware())

	// Health and metrics
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]any{"status": "ok"})
	})
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public
	auth := e.Group("/api/v1/auth")
	auth.POST("/login", d.AuthHandler.Login)
	auth.POST("/register", d.AuthHandler.Register)

	// Protected
	api := e.Group("/api/v1")
	api.Use(JWTAuth(d.JWT))

	api.GET("/me", d.AuthHandler.Me)
	api.GET("/user", d.AuthHandler.Me) // alias

	// Admin-only: create staff with role + location access
	api.POST("/auth/staff", d.AuthHandler.RegisterStaff, RequireRole("admin"))

	// Devices
	dev := api.Group("/devices")
	dev.GET("", d.DeviceHandler.List)
	dev.GET("/:id", d.DeviceHandler.Get)
	dev.POST("", d.DeviceHandler.Create)
	dev.PUT("/:id", d.DeviceHandler.Update)
	dev.DELETE("/:id", d.DeviceHandler.Delete)

	return e
}
