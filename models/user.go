package models

import (
	"time"
)

type UserLevel string

const (
	UserLevelAdmin   UserLevel = "ADMIN"
	UserLevelSupport UserLevel = "SUPPORT"
)

type UserRoleAction string

const (
	UserRoleAll    UserRoleAction = "ALL"
	UserRoleInsert UserRoleAction = "INSERT"
	UserRoleUpdate UserRoleAction = "UPDATE"
	UserRoleQuery  UserRoleAction = "QUERY"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusDisable UserStatus = "DISABLE"
	UserStatusSuspend UserStatus = "SUSPEND"
)

type User struct {
	UserID     int            `json:"user_id"`
	Username   string         `json:"username"`
	Telephone  string         `json:"telephone"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	Level      UserLevel      `json:"level,omitempty"`
	RoleAction UserRoleAction `json:"role_action,omitempty"`
	Status     UserStatus     `json:"status,omitempty"`
	OpID       *int           `json:"op_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
	DeletedAt  *time.Time     `json:"deleted_at,omitempty"`
}

type JwtUser struct {
	UserID     int            `json:"user_id"`
	Username   string         `json:"username"`
	Telephone  string         `json:"telephone"`
	Email      string         `json:"email"`
	Level      UserLevel      `json:"level"`
	RoleAction UserRoleAction `json:"role_action"`
}
