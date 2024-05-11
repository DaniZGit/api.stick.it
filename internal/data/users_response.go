package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Tokens 		int							 `json:"tokens"`
	Token			string					 `json:"token,omitempty"`
	AvailableFreePacks int `json:"available_free_packs"`
	LastFreePackObtainDate pgtype.Timestamp `json:"last_free_pack_obtain_date"`
}

type UserResponse struct {
	User User `json:"user"`
}

func CastToUserResponse(userRow database.User, token string) UserResponse {
	return UserResponse{
		User: User{
				ID: userRow.ID,
				CreatedAt: userRow.CreatedAt,
				Username: userRow.Username,
				Email: userRow.Email,
				Tokens: int(userRow.Tokens),
				Token: token,
				AvailableFreePacks: int(userRow.AvailableFreePacks),
				LastFreePackObtainDate: userRow.LastFreePackObtainDate,
		},
	}
}