package data

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Tokens 		int							 `json:"tokens"`
	Token			string					 `json:"token,omitempty"`
	AvailableFreePacks int `json:"available_free_packs"`
	LastFreePackObtainDate pgtype.Timestamp `json:"last_free_pack_obtain_date"`
}