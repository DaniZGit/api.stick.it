package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Page struct {
	ID        uuid.NullUUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	SortOrder int32 `json:"sort_order"`
	File *File `json:"file"` 
	Stickers []Sticker `json:"stickers"`
}

type PageResponse struct {
	Page Page `json:"page"`
}

type PagesResponse struct {
	Pages []Page `json:"pages"`
}

func BuildPageResponse(pageRows interface{}, file *database.File) any {
	switch value := pageRows.(type) {
		case database.Page:
			return PageResponse{
				Page: Page{
					ID: uuid.NullUUID{UUID: value.ID, Valid: !value.ID.IsNil()},
					CreatedAt: value.CreatedAt,
					SortOrder: value.SortOrder,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
					Stickers: []Sticker{},
				},
			}
		case []database.GetPageRow:
			return castToPageResponse(value)
		case []database.GetPagesRow:
			return castToPagesResponse(value)
	}

	return AlbumResponse{}
}

func castToPageResponse(pageRows []database.GetPageRow) PageResponse {
	if pageRows == nil || len(pageRows) <= 0 {
		return  PageResponse{
			Page: Page{},
		}
	}

	firstPageRow := pageRows[0]
	page := Page{
		ID: uuid.NullUUID{UUID: firstPageRow.ID, Valid: !firstPageRow.ID.IsNil()},
		CreatedAt: firstPageRow.CreatedAt,
		SortOrder: firstPageRow.SortOrder,
	}

	// add file
	if !firstPageRow.PageFileID.UUID.IsNil() {
		page.File = &File{
			ID: firstPageRow.PageFileID,
			Name: firstPageRow.PageFileName.String,
			Url: assetmanager.GetPublicAssetsFileUrl(firstPageRow.PageFilePath.String, ""),
		}
	}

	// add stickers
	if !firstPageRow.StickerID.UUID.IsNil() {
		for _, pageRow := range pageRows {
			sticker := Sticker{
				ID: pageRow.StickerID,
				CreatedAt: pageRow.StickerCreatedAt,
				Title: pageRow.StickerTitle.String,
				Type: pageRow.StickerType.String,
				Top: pageRow.StickerTop,
				Left: pageRow.StickerLeft,
				Width: pageRow.StickerWidth,
				Height: pageRow.StickerHeight,
				Numerator: pageRow.StickerNumerator.Int32,
				Denominator: pageRow.StickerDenominator.Int32,
				Rotation: pageRow.StickerRotation,
				File: &File{
					ID: pageRow.StickerFileID,
					Name: pageRow.StickerFileName.String,
					Url: assetmanager.GetPublicAssetsFileUrl(pageRow.StickerFilePath.String, ""),
				},
				PageID: pageRow.ID,
				RarityID: pageRow.StickerRarityID.UUID,
			}

			page.Stickers = append(page.Stickers, sticker)
		}
	} else {
		page.Stickers = []Sticker{}
	}

	return PageResponse{
		Page: page,
	}
}

func castToPagesResponse(pagesRows []database.GetPagesRow) PagesResponse {
	if pagesRows == nil || len(pagesRows) <= 0 {
		return PagesResponse{
			Pages: []Page{},
		}
	}

	pages := []Page{}
	for _, pagesRow := range pagesRows {
		page := Page{
			ID: uuid.NullUUID{UUID: pagesRow.ID, Valid: !pagesRow.ID.IsNil()},
			CreatedAt: pagesRow.CreatedAt,
			SortOrder: pagesRow.SortOrder,
		}

		// add file
		if !pagesRow.PageFileID.UUID.IsNil() {
			page.File = &File{
				ID: pagesRow.PageFileID,
				Name: pagesRow.PageFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(pagesRow.PageFilePath.String, ""),
			}
		}

		pages = append(pages, page)
	}
	
	return PagesResponse{
		Pages: pages,
	}
}
