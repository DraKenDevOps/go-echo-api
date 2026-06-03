package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"go-echo-api/models"
)

func (h *Handler) SaveNewAccount(c echo.Context) error {
	var ac models.Account
	ctx := c.Request().Context()
	err := c.Bind(&ac)
	if err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Invalid request body",
		})
	}

	if h.cfg.LimitMaxBalance && ac.Balance > 10000 {
		h.logger.Warn("Balance limit on account creating")
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Balance over the limit",
		})
	}

	const sql = "INSERT INTO accounts (balance) VALUES (?)"

	result, err := h.db.ExecContext(ctx, sql, ac.Balance)
	if err != nil {
		h.logger.Error("Error save new account", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status":  "error",
			"message": "Failed to save account",
		})
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		h.logger.Error("Error getting insert id", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status":  "error",
			"message": "Failed to get insert id",
		})
	}

	h.logger.Info("Save new account", zap.Int64("accountId", insertId))

	return c.JSON(200, map[string]any{
		"status": "success", "message": "Save new successfully",
	})
}
