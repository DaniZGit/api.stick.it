package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Avatar struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title  string           `json:"title"`
	File *File `json:"file"`
}

type AvatarsResponse struct {
	Avatars []Avatar `json:"avatars"`
	Metadata Metadata `json:"metadata"`
}

type AvatarResponse struct {
	Avatar Avatar `json:"avatar"`
}

func CastToAvatarResponse(row database.Avatar, file database.File) AvatarResponse {
	avatar := Avatar{
		ID: row.ID,
		CreatedAt: row.CreatedAt,
		Title: row.Title,
	}

	if !file.ID.IsNil() {
		avatar.File = &File{
			ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
			Name: file.Name,
			Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),	
		}
	}

	return AvatarResponse{
		Avatar: avatar,
	}
}

func CastToAvatarsResponse(rows []database.GetAvatarsRow, metadata Metadata) AvatarsResponse {
	if rows == nil || len(rows) <= 0 {
		return  AvatarsResponse{
			Metadata: metadata,
			Avatars: []Avatar{},
		}
	}

	avatars := []Avatar{}
	for _, row := range rows {
		avatar := Avatar{
			ID: row.ID,
			CreatedAt: row.CreatedAt,
			Title: row.Title,
			File: &File{
				ID: row.FileID,
				Name: row.AvatarFileName,
				Url: assetmanager.GetPublicAssetsFileUrl(row.AvatarFilePath, ""),
			},
		}

		avatars = append(avatars, avatar)
	}

	return AvatarsResponse{
		Metadata: metadata,
		Avatars: avatars,
	}
}