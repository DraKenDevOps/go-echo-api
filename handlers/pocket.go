package handlers

import (
	"go-echo-api/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) SaveNewPocket(c echo.Context) error {
	var pk models.Pocket
	ctx := c.Request().Context()
	err := c.Bind(&pk)
	if err != nil {
		h.Logger.Error("Invalid request body", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Invalid request body",
		})
	}

	if pk.Currency == "" {
		pk.Currency = "USD"
	}

	if pk.Amount > 0.0 {
		pk.Amount = 0.0
	}

	if len(pk.Currency) != 3 {
		h.Logger.Error("Invlid request body currency")
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Currency must be 3 characters",
		})
	}

	var insertId int

	const sql = "INSERT INTO pockets(pocket_name, amount, currency, account_id) VALUES (?,?,?,?)"
	err = h.DB.QueryRowContext(ctx, sql, pk.PocketName, pk.Amount, pk.Currency, pk.AccountId).Scan(&insertId)

	if err != nil {
		h.Logger.Error("Error save new pocket", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Failed to save pocket",
		})
	}

	h.Logger.Info("Save new account", zap.Int("pocketId", insertId))

	return c.JSON(200, map[string]any{
		"status": "success", "message": "Save new successfully",
	})

}
