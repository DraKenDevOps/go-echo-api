package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go-echo-api/config"
	"go-echo-api/database"
	"go-echo-api/routes"
	"go-echo-api/utils"
	"go-echo-api/zplogger"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()
	logger := zplogger.NewLogger(cfg)
	defer logger.Close()

	db := database.InitDB(cfg, logger)
	defer db.Close()

	e := echo.New()

	e.Static("/static", filepath.Join(cfg.Cwd, "uploads"))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			return nil
		},
	}))
	e.Use(routes.LoggerMiddleware(logger))

	startTime := time.Now()
	e.GET("/health", func(c echo.Context) error {
		res := map[string]any{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"uptime":    utils.FormatUptime(time.Since(startTime)),
			"version":   cfg.AppVersion,
		}
		return c.JSON(200, res)
	})

	routes.RegisterRoutes(e, db, cfg, logger)

	hostAddr := cfg.Host + ":" + cfg.ServerPort
	logger.Info("Server starting on " + hostAddr)

	err := e.Start(hostAddr)
	go func() {
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server startup failed", zap.Error(err))
		}
	}()

	// ── Graceful shutdown ──────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Warn("Received signal, shutting down", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}
}
