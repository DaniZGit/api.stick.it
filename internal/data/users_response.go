package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Description pgtype.Text `json:"description"`
	Tokens 		int							 `json:"tokens"`
	Token			string					 `json:"token,omitempty"`
	AvailableFreePacks int `json:"available_free_packs"`
	LastFreePackObtainDate pgtype.Timestamp `json:"last_free_pack_obtain_date"`
	AvatarID uuid.UUID `json:"avatar_id"`
	Avatar Avatar `json:"avatar"`
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
			Description: userRow.Description,
			Tokens: int(userRow.Tokens),
			Token: token,
			AvailableFreePacks: int(userRow.AvailableFreePacks),
			LastFreePackObtainDate: userRow.LastFreePackObtainDate,
			AvatarID: userRow.AvatarID.UUID,
		},
	}
}

func CastToUserByIDResponse(userRow database.GetUserByIDRow, token string) UserResponse {
	return UserResponse{
		User: User{
			ID: userRow.ID,
			CreatedAt: userRow.CreatedAt,
			Username: userRow.Username,
			Email: userRow.Email,
			Description: userRow.Description,
			Tokens: int(userRow.Tokens),
			Token: token,
			AvailableFreePacks: int(userRow.AvailableFreePacks),
			LastFreePackObtainDate: userRow.LastFreePackObtainDate,
			AvatarID: userRow.AvatarID.UUID,
			Avatar: Avatar{
				ID: userRow.AvatarID.UUID,
				Title: userRow.AvatarTitle.String,
				File: &File{
					ID: userRow.AvatarFileID,
					Name: userRow.AvatarFileName.String,
					Url: assetmanager.GetPublicAssetsFileUrl(userRow.AvatarFilePath.String, ""),
				},
			},
		},
	}
}