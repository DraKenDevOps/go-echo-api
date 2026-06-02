package handlers

import (
	"database/sql"

	"go-echo-api/config"
	"go-echo-api/zplogger"
)

// Handler holds shared dependencies used by all sub-handlers.
type Handler struct {
	DB     *sql.DB
	Config *config.Config
	Logger *zplogger.Logger
}
