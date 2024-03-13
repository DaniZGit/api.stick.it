package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Album struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title  string           `json:"title"`
	DateFrom string `json:"date_from"`
	DateTo string `json:"date_to"`
	Featured bool `json:"featured"`
	File *File `json:"file"`
	Pages []Page `json:"pages"`
}

type AlbumsResponse struct {
	Albums []Album `json:"albums"`
}

type AlbumResponse struct {
	Album Album `json:"album"`
}

func BuildAlbumResponse(albumRows interface{}, file *database.File) any {
	switch value := albumRows.(type) {
		case database.Album:
			return AlbumResponse{
				Album: Album{
					ID: value.ID,
					CreatedAt: value.CreatedAt,
					Title: value.Title,
					DateFrom: value.DateFrom.Time.String(),
					DateTo: value.DateTo.Time.String(),
					Featured: value.Featured.Bool,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
					Pages: []Page{},
				},
			}
		case []database.GetAlbumRow:
			return castToAlbumResponse(value)
		case []database.GetAlbumsRow:
			return castToAlbumsResponse(value)
	}

	return AlbumResponse{}
}

func castToAlbumResponse(albumRows []database.GetAlbumRow) AlbumResponse {
	if albumRows == nil || len(albumRows) <= 0 {
		return  AlbumResponse{
			Album: Album{},
		}
	}

	album := Album{
		ID: albumRows[0].ID,
		CreatedAt: albumRows[0].CreatedAt,
		Title: albumRows[0].Title,
		DateFrom: albumRows[0].DateFrom.Time.String(),
		DateTo: albumRows[0].DateTo.Time.String(),
		Featured: albumRows[0].Featured.Bool,
	}

	// add file
	if !albumRows[0].AlbumFileID.UUID.IsNil() {
		album.File = &File{
			ID: albumRows[0].AlbumFileID,
			Name: albumRows[0].AlbumFileName.String,
			Url: assetmanager.GetPublicAssetsFileUrl(albumRows[0].AlbumFilePath.String, ""),
		}
	}

	// add pages
	if !albumRows[0].PageID.UUID.IsNil() {
		for _, albumRow := range albumRows {
			page := Page{
				ID: albumRow.PageID,
				CreatedAt: albumRow.PageCreatedAt,
				SortOrder: albumRow.PageSortOrder.Int32,
				File: &File{
					ID: albumRow.PageFileID,
					Name: albumRow.PageFileName.String,
					Url: assetmanager.GetPublicAssetsFileUrl(albumRow.PageFilePath.String, ""),
				},
			}

			album.Pages = append(album.Pages, page)
		}
	} else {
		album.Pages = []Page{}
	}

	return AlbumResponse{
		Album: album,
	}
}

func castToAlbumsResponse(albumsRows []database.GetAlbumsRow) AlbumsResponse {
	if albumsRows == nil || len(albumsRows) <= 0 {
		return AlbumsResponse{
			Albums: []Album{},
		}
	}

	albums := []Album{}
	for _, albumsRow := range albumsRows {
		album := Album{
			ID: albumsRow.ID,
			CreatedAt: albumsRow.CreatedAt,
			Title: albumsRow.Title,
			DateFrom: albumsRow.DateFrom.Time.String(),
			DateTo: albumsRow.DateTo.Time.String(),
			Featured: albumsRow.Featured.Bool,
		}

		// add file
		if !albumsRow.AlbumFileID.UUID.IsNil() {
			album.File = &File{
				ID: albumsRow.AlbumFileID,
				Name: albumsRow.AlbumFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(albumsRow.AlbumFilePath.String, ""),
			}
		}

		albums = append(albums, album)
	}
	
	return AlbumsResponse{
		Albums: albums,
	}
}
