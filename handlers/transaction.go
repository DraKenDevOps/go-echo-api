package handlers

import (
	"database/sql"
	"fmt"
	"go-echo-api/models"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type TxnResponse struct {
	Data      []models.Transaction `json:"data"`
	TotalPage int                  `json:"totalPage"`
}

const (
	sqlTxn   = "INSERT INTO transactions (from_pocket_id, to_pocket_id, amount) VALUES (?, ?, ?)"
	listStmt = `SELECT txn_id, from_pocket_id, to_pocket_id, amount, DATE_FORMAT(txn_date, '%Y-%m-%d %H:%i:%s')
		FROM transactions
		WHERE (from_pocket_id ? OR to_pocket_id = ?)
		ORDER BY id DESC
		LIMIT ? OFFSET ?`

	countStmt = "SELECT COUNT(*) FROM transactions WHERE (from_pocket_id ? OR to_pocket_id = ?)"
)

func (h *Handler) SaveTransaction(ec echo.Context) error {
	ctx := ec.Request().Context()

	var txn models.Transaction
	err := ec.Bind(&txn)

	if err != nil {
		h.logger.Error("bad request body", zap.Error(err))
		return ec.JSON(200, map[string]any{"status": "error", "message": "bad request body"})
	}

	if txn.Amount <= 0 {
		h.logger.Error("amount must more than 0")
		return ec.JSON(200, map[string]any{"status": "error", "message": "bad request body"})
	}

	if txn.ToPocketID <= 0 {
		h.logger.Error("invalid pocket id")
		return ec.JSON(200, map[string]any{"status": "error", "message": "bad request body"})
	}

	_, errx := h.db.ExecContext(ctx, sqlTxn, txn.FromPocketID, txn.ToPocketID, txn.Amount)
	if errx != nil {
		h.logger.Error("Cannot insert transaction", zap.Error(errx))
		return ec.JSON(200, map[string]any{"status": "error", "message": "Cannot insert transaction"})
	}

	fp, err := getPocketById(h.db, txn.FromPocketID)
	if err != nil {
		h.logger.Error("Internal error get pocket FromPocketID", zap.Error(err))
		return ec.JSON(
			200,
			map[string]any{"status": "error", "message": err.Error()},
		)
	}

	tp, err := getPocketById(h.db, txn.ToPocketID)
	if err != nil {
		h.logger.Error("Internal error get pocket ToPocketID", zap.Error(err))
		return ec.JSON(
			200,
			map[string]any{"status": "error", "message": err.Error()},
		)
	}

	sqlTx, err := h.db.Begin()
	if err != nil {
		h.logger.Error("Cannot begin sql transactions", zap.Error(err))
		return ec.JSON(200, map[string]any{"status": "error", "message": "Internal server error"})
	}
	txAmount := decimal.NewFromFloat(txn.Amount)
	fromAmount := decimal.NewFromFloat(fp.Amount)
	fromAmount = fromAmount.Sub(txAmount)

	updateAmountPocketById(sqlTx, fp.PocketId, fromAmount.InexactFloat64())

	toAmount := decimal.NewFromFloat(tp.Amount)
	toAmount = decimal.Sum(toAmount, txAmount)
	fmt.Println(toAmount.InexactFloat64())

	updateAmountPocketById(sqlTx, tp.PocketId, toAmount.InexactFloat64())

	if err = sqlTx.Commit(); err != nil {
		h.logger.Error("Cannot COMMIT sql transactions", zap.Error(err))
		return ec.JSON(200, map[string]any{"status": "error", "message": "Internal server error"})
	}

	return ec.JSON(200, map[string]any{"status": "success", "message": "Save transaction successfully"})
}

func updateAmountPocketById(tx *sql.Tx, pocketId int, amount float64) {
	res, err := tx.Exec("UPDATE pockets SET amount = ? WHERE id = ?", amount, pocketId)
	if err != nil {
		fmt.Println(err.Error())
		tx.Rollback()
	}
	_, errrf := res.RowsAffected()
	if errrf != nil {
		fmt.Println(err.Error())
		tx.Rollback()
	}
}

func getPocketById(db *sql.DB, id int) (*models.Pocket, error) {
	var pk models.Pocket
	err := db.QueryRow(sqlId, id).Scan(&pk.PocketId, &pk.PocketName, &pk.Amount, &pk.AccountId, &pk.Currency)
	if err != nil {
		return nil, err
	}
	return &pk, nil
}

func (h *Handler) GetTransactionList(ec echo.Context) error {
	ctx := ec.Request().Context()
	pocketId := ec.Param("id")

	pageStr := ec.QueryParam("page")

	page := 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	offset := (page - 1) * limit
	var txns []models.Transaction

	var total int
	err := h.db.QueryRowContext(
		ctx,
		countStmt,
		pocketId,
		pocketId,
	).Scan(&total)

	if err != nil {
		h.logger.Error("count transactions error", zap.Error(err))
		return ec.JSON(200, map[string]any{"data": txns, "totalPage": 1})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	rows, err := h.db.QueryContext(
		ctx,
		listStmt,
		pocketId,
		pocketId,
		limit,
		offset,
	)

	if err != nil {
		h.logger.Error("query transactions error", zap.Error(err))
		return ec.JSON(200, map[string]any{"data": txns, "totalPage": 1})
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
			return ec.JSON(200, map[string]any{"data": txns, "totalPage": 1})
		}
	}

	return ec.JSON(200, TxnResponse{
		Data:      txns,
		TotalPage: totalPage,
	})
}
