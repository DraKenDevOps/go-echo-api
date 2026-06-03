package handlers

import (
	"database/sql"

	"go-echo-api/config"
	"go-echo-api/zplogger"
)

// Handler holds shared dependencies used by all sub-handlers.
type Handler struct {
	db     *sql.DB
	cfg    *config.Config
	logger *zplogger.Logger
}

func NewHandler(db *sql.DB, cfg *config.Config, logger *zplogger.Logger) *Handler {
	return &Handler{db: db, cfg: cfg, logger: logger}
}
