package data

import (
	"time"

	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AlbumResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title  string           `json:"title"`
	DateFrom string `json:"date_from"`
	DateTo string `json:"date_to"`
	Featured bool `json:"featured"`
	File FileResponse `json:"file"` 
}

func CastToAlbumResponse(album database.GetAlbumRow) AlbumResponse {
	return AlbumResponse{
		ID: album.ID,
		CreatedAt: album.CreatedAt,
		DateFrom: album.DateFrom.Time.Format(time.RFC3339),
		DateTo: album.DateTo.Time.Format(time.RFC3339),
		Featured: album.Featured.Bool,
		Title: album.Title,
		File: FileResponse{
			ID: album.Albumfile.ID,
			Name: album.Albumfile.Name.String,
			Url: assetmanager.GetPublicAssetsFileUrl(album.Albumfile.Path.String, true),
		},
	}
}

func CastToAlbumsResponse(albums []database.GetAlbumsRow) []AlbumResponse {
	castedAlbums := make([]AlbumResponse, 0)
	for _, album := range albums {
		a := AlbumResponse{
			ID: album.ID,
			CreatedAt: album.CreatedAt,
			DateFrom: album.DateFrom.Time.String(),
			DateTo: album.DateTo.Time.String(),
			Featured: album.Featured.Bool,
			Title: album.Title,
			File: FileResponse{
				ID: album.Albumfile.ID,
				Name: album.Albumfile.Name.String,
				Url: assetmanager.GetPublicAssetsFileUrl(album.Albumfile.Path.String, true),
			},
		}

		castedAlbums = append(castedAlbums, a)
	}

	return castedAlbums
}

func CastToAlbumUpdateResponse(album database.Album, file database.File) AlbumResponse {
	return AlbumResponse{
		ID: album.ID,
		CreatedAt: album.CreatedAt,
		DateFrom: album.DateFrom.Time.Format(time.RFC3339),
		DateTo: album.DateTo.Time.Format(time.RFC3339),
		Featured: album.Featured.Bool,
		Title: album.Title,
		File: FileResponse{
			ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
			Name: file.Name,
			Url: assetmanager.GetPublicAssetsFileUrl(file.Path, true),
		},
	}
}