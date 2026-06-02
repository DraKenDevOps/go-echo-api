package models

import "time"

type Transaction struct {
	TxnID        int       `json:"txn_id"`
	FromPocketID int       `json:"from_pocket_id"`
	ToPocketID   int       `json:"to_pocket_id"`
	Amount       float64   `json:"amount"`
	TxnDate      time.Time `json:"txn_date"`
}
