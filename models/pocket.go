package models

import "time"

type Pocket struct {
	PocketId   int        `json:"pocket_id"`
	PocketName string     `json:"pocket_name"`
	Amount     float64    `json:"amount"`
	Currency   string     `json:"currency"`
	AccountId  *int       `json:"account_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
