package routes

import (
	"go-echo-api/config"
	"go-echo-api/handlers"
	middleware "go-echo-api/middlewares"
	"go-echo-api/zplogger"

	"database/sql"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(e *echo.Echo, db *sql.DB, logger *zplogger.Logger, cfg *config.Config) {
	handler := &handlers.Handler{DB: db, Config: cfg, Logger: logger}

	api := e.Group(cfg.BasePath)
	api.POST("/login", handler.Login)
	api.GET("/refresh", handler.Refresh, middleware.AuthChecker(cfg, logger))

	// TODO: Add other route groups (accounts, pockets, transactions)
}

// LoggerMiddleware returns a middleware that logs requests using zap logger
// with the provided configuration.
func LoggerMiddleware(logger *zplogger.Logger) echo.MiddlewareFunc {
	return middleware.LoggerMiddleware(middleware.LoggerConfig{
		Logger:         logger,
		IgnorePaths:    []string{"/health", "/metrics"},
		IgnoreBodyKeys: []string{"file", "files"},
		MaskedKeys:     []string{"password"},
	})
}
