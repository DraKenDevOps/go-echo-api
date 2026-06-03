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
func RegisterRoutes(e *echo.Echo, db *sql.DB, cfg *config.Config, logger *zplogger.Logger) {
	handler := handlers.NewHandler(db, cfg, logger)

	api := e.Group(cfg.BasePath)
	api.POST("/login", handler.Login)
	api.GET("/refresh", handler.Refresh, middleware.AuthChecker(cfg, logger))

	api.POST("/save_account", handler.SaveNewAccount, middleware.AuthChecker(cfg, logger))
	api.POST("/save_pocket", handler.SaveNewPocket, middleware.AuthChecker(cfg, logger))
	api.POST("/save_transaction", handler.SaveTransaction, middleware.AuthChecker(cfg, logger))

	api.GET("/get_pocket/:id", handler.GetPocketById, middleware.AuthChecker(cfg, logger))
	api.GET("/get_transactions/:id", handler.GetTransactionList, middleware.AuthChecker(cfg, logger))
}

// LoggerMiddleware returns a middleware that logs requests using zap logger
// with the provided configuration.
func LoggerMiddleware(logger *zplogger.Logger) echo.MiddlewareFunc {
	return middleware.LoggerMiddleware(middleware.LoggerConfig{
		Logger:         logger,
		IgnorePaths:    []string{"/health", "/metrics"},
		IgnoreBodyKeys: []string{"file", "files"},
		MaskedKeys:     []string{"password", "token", "access_token", "accessToken"},
	})
}
