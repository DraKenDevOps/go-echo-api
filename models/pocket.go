package models

import "time"

type Pocket struct {
	PocketId   int        `json:"pocket_id,omitempty"`
	PocketName string     `json:"pocket_name"`
	Amount     float64    `json:"amount"`
	Currency   string     `json:"currency"`
	AccountId  *int       `json:"account_id"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type PocketReqBody struct {
	PocketId   int     `json:"pocket_id,omitempty"`
	PocketName string  `json:"pocket_name"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	AccountId  *int    `json:"account_id"`
}
