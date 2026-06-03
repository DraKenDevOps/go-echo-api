package handlers

import (
	"go-echo-api/models"
	"math"
	"strconv"

	// "regexp"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cntStmt = "SELECT COUNT(*) FROM pockets"
	lstStmt = "SELECT pocket_id, pocket_name, amount, account_id, currency FROM pockets LIMIT ? OFFSET ?"
	sqlId   = "SELECT pocket_id, pocket_name, amount, account_id, currency FROM pockets WHERE pocket_id = ?"
)

type PkResponse struct {
	Data      []models.Pocket `json:"data"`
	TotalPage int             `json:"totalPage"`
}

func (h *Handler) GetPocketList(ec echo.Context) error {
	pageStr := ec.QueryParam("page")
	ctx := ec.Request().Context()

	page := 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	offset := (page - 1) * limit
	var pks []models.Pocket

	var total int
	err := h.db.QueryRowContext(
		ctx,
		countStmt,
	).Scan(&total)

	if err != nil {
		h.logger.Error("count transactions error", zap.Error(err))
		return ec.JSON(200, map[string]any{"data": pks, "totalPage": 1})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	rows, err := h.db.QueryContext(
		ctx,
		listStmt,
		limit,
		offset,
	)

	if err != nil {
		h.logger.Error("query transactions error", zap.Error(err))
		return ec.JSON(200, map[string]any{"data": pks, "totalPage": 1})
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		var gotDate string

		err := rows.Scan(
			&t.TxnID,
			&t.FromPocketID,
			&t.ToPocketID,
			&t.Amount,
			&gotDate,
		)

		if err != nil {
			h.logger.Error("scan transaction error", zap.Error(err))
			return ec.JSON(200, map[string]any{"data": pks, "totalPage": 1})
		}
	}

	return ec.JSON(200, PkResponse{
		Data:      pks,
		TotalPage: totalPage,
	})
}

func (h *Handler) GetPocketById(ec echo.Context) error {
	ctx := ec.Request().Context()
	id := ec.Param("id")
	var pk models.Pocket

	err := h.db.QueryRowContext(ctx, sqlId, id).Scan(&pk.PocketId, &pk.PocketName, &pk.Amount, &pk.AccountId, &pk.Currency)
	if err != nil {
		// match, errMatch := regexp.MatchString("invalid input syntax", err.Error())
		// if match {
		// 	h.logger.Error("Param id must be integer", zap.Error(err))
		// 	return ec.JSON(
		// 		200,
		// 		map[string]any{"status": "error", "message": "Param id must be integer"},
		// 	)
		// }
		// if errMatch != nil {
		// 	h.logger.Error("Match fail", zap.Error(err))
		// 	return ec.JSON(
		// 		200,
		// 		map[string]any{"status": "error", "message": err.Error()},
		// 	)
		// }
		// match, errMatch = regexp.MatchString("no rows in result set", err.Error())
		// if match {
		// 	h.logger.Error("Record not found", zap.Error(err))
		// 	return ec.JSON(
		// 		200,
		// 		map[string]any{"status": "error", "message": "Record not found"},
		// 	)
		// }
		// if errMatch != nil {
		// 	h.logger.Error("Match fail", zap.Error(err))
		// 	return ec.JSON(
		// 		200,
		// 		map[string]any{"status": "error", "message": err.Error()},
		// 	)
		// }
		h.logger.Error("Internal error", zap.Error(err))
		return ec.JSON(
			200,
			map[string]any{"status": "error", "message": err.Error()},
		)
	}

	return ec.JSON(200, pk)
}

func (h *Handler) SaveNewPocket(c echo.Context) error {
	var pk models.Pocket
	ctx := c.Request().Context()
	err := c.Bind(&pk)
	if err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
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
		h.logger.Error("Invlid request body currency")
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Currency must be 3 characters",
		})
	}

	const sql = "INSERT INTO pockets(pocket_name, amount, currency, account_id) VALUES (?,?,?,?)"
	result, err := h.db.ExecContext(ctx, sql, pk.PocketName, pk.Amount, pk.Currency, pk.AccountId)

	if err != nil {
		h.logger.Error("Error save new pocket", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status": "error", "message": "Failed to save pocket",
		})
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		h.logger.Error("Error getting insert id", zap.Error(err))
		return c.JSON(200, map[string]any{
			"status":  "error",
			"message": "Failed to save pocket",
		})
	}

	h.logger.Info("Save new account", zap.Int64("pocketId", insertId))

	return c.JSON(200, map[string]any{
		"status": "success", "message": "Save new successfully",
	})
}
