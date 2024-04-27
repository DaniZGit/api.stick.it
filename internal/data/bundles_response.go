package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Bundle struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title string `json:"title"`
	Price pgtype.Numeric `json:"price"`
	Tokens int `json:"tokens"`
	Bonus int `json:"bonus"`
	FileID uuid.NullUUID `json:"file_id"`
	File *File `json:"file"`
}

type BundleResponse struct {
	Bundle Bundle `json:"bundle"`
}

type BundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
	Metadata Metadata `json:"metadata"`
}

func BuildBundlesResponse(rows interface{}, metadata Metadata, file *database.File) any {
	switch value := rows.(type) {
		case database.Bundle:
			return BundleResponse{
				Bundle: Bundle{
					ID: value.ID,
					CreatedAt: value.CreatedAt,
					Title: value.Title,
					Price: value.Price,
					Tokens: int(value.Tokens),
					Bonus: int(value.Bonus),
					FileID: value.FileID,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
				},
			}
		case []database.GetBundlesRow:
			return castToBundlesResponse(value, metadata)
	}

	return BundlesResponse{}
}

func castToBundlesResponse(rows []database.GetBundlesRow, metadata Metadata) BundlesResponse {
	if rows == nil || len(rows) <= 0 {
		return BundlesResponse{
			Bundles: []Bundle{},
		}
	}

	bundles := []Bundle{}
	for _, row := range rows {
		bundle := Bundle{
				ID: row.ID,
				CreatedAt: row.CreatedAt,
				Title: row.Title,
				Price: row.Price,
				Tokens: int(row.Tokens),
				Bonus: int(row.Bonus),
				FileID: uuid.NullUUID{UUID: row.FileID.UUID, Valid: !row.FileID.UUID.IsNil()},
		}

		// add file
		if !row.BundleFileID.UUID.IsNil() {
			bundle.File = &File{
				ID: row.BundleFileID,
				Name: row.BundleFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(row.BundleFilePath.String, ""),
			}
		}

		bundles = append(bundles, bundle)
	}

	return BundlesResponse{
		Bundles: bundles,
		Metadata: metadata,
	}
}
