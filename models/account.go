package models

import "time"

type Account struct {
	AccountId int        `json:"account_id,omitempty"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"created_at,omitzero"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AccountReqBody struct {
	AccountId int     `json:"account_id,omitempty"`
	Balance   float64 `json:"balance"`
}
