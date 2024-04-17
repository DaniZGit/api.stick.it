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
	Width pgtype.Numeric `json:"width"`
	Height pgtype.Numeric `json:"height"`
	Numerator int32 `json:"numerator"`
	Denominator int32 `json:"denominator"`
	Rotation pgtype.Numeric `json:"rotation"`
	PageID uuid.UUID `json:"page_id"`
	RarityID uuid.NullUUID `json:"rarity_id"`
	StickerID uuid.NullUUID `json:"sticker_id"`
	FileID uuid.NullUUID `json:"file_id"`
	File *File `json:"file"`
	Rarity *Rarity `json:"rarity"`
}

type UserSticker struct {
	ID        uuid.UUID        `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	StickerID uuid.UUID `json:"sticker_id"`
	Amount int `json:"amount"`
	Sticker Sticker `json:"sticker"`
}

type StickerResponse struct {
	Sticker Sticker `json:"sticker"`
}

type StickersResponse struct {
	Stickers []Sticker `json:"stickers"`
}

type UserStickersResponse struct {
	UserStickers []UserSticker `json:"user_stickers"`
}

func BuildStickerResponse(stickerRows interface{}, file *database.File, rarity *database.Rarity) any {
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
					Width: value.Width,
					Height: value.Height,
					Numerator: value.Numerator,
					Denominator: value.Denominator,
					Rotation: value.Rotation,
					PageID: value.PageID,
					RarityID: value.RarityID,
					FileID: value.FileID,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
					Rarity: &Rarity{
						ID: uuid.NullUUID{UUID: rarity.ID, Valid: !rarity.ID.IsNil()},
						Title: rarity.Title,
					},
				},
			}
		case []database.GetPageStickersRow:
			return castToStickersResponse(value)
		case []database.GetStickerRaritiesRow:
			return castToStickerRaritiesResponse(value)
		case []database.GetUserStickersRow:
			return castToUserStickersResponse(value)
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
			Width: stickersRow.Width,
			Height: stickersRow.Height,
			Numerator: stickersRow.Numerator,
			Denominator: stickersRow.Denominator,
			Rotation: stickersRow.Rotation,
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

func castToStickerRaritiesResponse(stickersRows []database.GetStickerRaritiesRow) StickersResponse {
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
			Width: stickersRow.Width,
			Height: stickersRow.Height,
			Numerator: stickersRow.Numerator,
			Denominator: stickersRow.Denominator,
			Rotation: stickersRow.Rotation,
			PageID: stickersRow.PageID,
			RarityID: stickersRow.RarityID,
			CreatedAt: stickersRow.CreatedAt,
			StickerID: stickersRow.StickerID,
		}

		// add rarity
		if !stickersRow.StickerRarityID.UUID.IsNil() {
			sticker.Rarity = &Rarity{
				ID: stickersRow.StickerRarityID,
				Title: stickersRow.StickerRarityTitle.String,
			}
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

func castToUserStickersResponse(rows []database.GetUserStickersRow) UserStickersResponse {
	if rows == nil || len(rows) <= 0 {
		return UserStickersResponse{
			UserStickers: []UserSticker{},
		}
	}

	userStickers := []UserSticker{}
	for _, row := range rows {
		sticker := Sticker{
			ID: uuid.NullUUID{UUID: row.ID, Valid: !row.ID.IsNil()},
			Title: row.StickerTitle,
			Type: row.StickerType,
			Top: row.StickerTop,
			Left: row.StickerLeft,
			Width: row.StickerWidth,
			Height: row.StickerHeight,
			Numerator: row.StickerNumerator,
			Denominator: row.StickerDenominator,
			Rotation: row.StickerRotation,
			PageID: row.StickerPageID,
			RarityID: row.StickerRarityID,
			CreatedAt: row.StickerCreatedAt,
			StickerID: row.StickerStickerID,
		}

		// add rarity
		if !row.StickerRarityID.UUID.IsNil() {
			sticker.Rarity = &Rarity{
				ID: row.StickerRarityID,
				Title: row.StickerRarityTitle.String,
			}
		}

		// add file
		if !row.StickerFileID.UUID.IsNil() {
			sticker.File = &File{
				ID: row.StickerFileID,
				Name: row.StickerFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(row.StickerFilePath.String, ""),
			}
		}

		userSticker := UserSticker{
			ID: row.ID,
			UserID: row.UserID,
			StickerID: row.StickerID,
			Amount: int(row.Amount),
			Sticker: sticker,
		}

		userStickers = append(userStickers, userSticker)
	}
	
	return UserStickersResponse{
		UserStickers: userStickers,
	}
}