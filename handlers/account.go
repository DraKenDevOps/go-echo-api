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
		h.Logger.Error("Invalid request body", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Invalid request body",
		})
	}

	if h.Config.LimitMaxBalance && ac.Balance > 10000 {
		h.Logger.Warn("Balance limit on account creating")
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Balance over the limit",
		})
	}

	const sql = "INSERT INTO accounts (balance) VALUES (?);"

	var insertId int
	err = h.DB.QueryRowContext(ctx, sql, ac.Balance).Scan(&insertId)
	if err != nil {
		h.Logger.Error("Error save new account", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Failed to save account",
		})
	}

	h.Logger.Info("Save new account", zap.Int("accountId", insertId))

	return c.JSON(200, map[string]any{
		"status": "success", "message": "Save new successfully",
	})
}
