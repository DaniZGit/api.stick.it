package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Sticker struct {
	ID        uuid.NullUUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title  string           `json:"title"`
	Type string `json:"type"`
	Top pgtype.Numeric `json:"top"`
	Left pgtype.Numeric `json:"left"`
	File *File `json:"file"`
	PageID uuid.UUID `json:"page_id"`
	RarityID uuid.UUID `json:"rarity_id"`
}

type StickerResponse struct {
	Sticker Sticker `json:"sticker"`
}

type StickersResponse struct {
	Stickers []Sticker `json:"stickers"`
}

func BuildStickerResponse(stickerRows interface{}, file *database.File) any {
	switch value := stickerRows.(type) {
		case database.Sticker:
			return StickerResponse{
				Sticker: Sticker{
					ID: uuid.NullUUID{UUID: value.ID, Valid: !value.ID.IsNil()},
					CreatedAt: value.CreatedAt,
					Title: value.Title,
					Type: value.Type,
					Top: value.Top,
					Left: value.Left,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
					PageID: value.PageID,
					RarityID: value.RarityID,
				},
			}
		case []database.GetPageStickersRow:
			return castToStickersResponse(value)
	}

	return StickerResponse{}
}

func castToStickersResponse(stickersRows []database.GetPageStickersRow) StickersResponse {
	if stickersRows == nil || len(stickersRows) <= 0 {
		return StickersResponse{
			Stickers: []Sticker{},
		}
	}

	stickers := []Sticker{}
	for _, stickersRow := range stickersRows {
		sticker := Sticker{
			ID: uuid.NullUUID{UUID: stickersRow.ID, Valid: !stickersRow.ID.IsNil()},
			Title: stickersRow.Title,
			Type: stickersRow.Type,
			Top: stickersRow.Top,
			Left: stickersRow.Left,
			PageID: stickersRow.PageID,
			RarityID: stickersRow.RarityID,
			CreatedAt: stickersRow.CreatedAt,
		}

		// add file
		if !stickersRow.StickerFileID.UUID.IsNil() {
			sticker.File = &File{
				ID: stickersRow.StickerFileID,
				Name: stickersRow.StickerFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(stickersRow.StickerFilePath.String, ""),
			}
		}

		stickers = append(stickers, sticker)
	}
	
	return StickersResponse{
		Stickers: stickers,
	}
}
