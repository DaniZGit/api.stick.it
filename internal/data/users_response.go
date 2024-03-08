package data

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Token			string					 `json:"token,omitempty"`
}