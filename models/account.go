package models

import "time"

type Account struct {
	AccountId int        `json:"account_id"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
