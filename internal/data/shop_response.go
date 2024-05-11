package data

import (
	"encoding/json"

	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

type ShopPacksResponse struct {
	Packs []Pack `json:"packs"`
}

type ShopBundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
}

func BuildShopResponse(rows interface{}) any {
	switch value := rows.(type) {
		case []database.GetShopPacksRow:
			return castToShopPacksResponse(value)
		case []database.GetShopBundlesRow:
			return castToShopBundlesResponse(value)
	}

	return ShopPacksResponse{}
}

func castToShopPacksResponse(rows []database.GetShopPacksRow) ShopPacksResponse {
	if rows == nil || len(rows) <= 0 {
		return ShopPacksResponse{
			Packs: []Pack{},
		}
	}

	packs := []Pack{}
	for _, row := range rows {
		pack := Pack{
				ID: row.ID,
				CreatedAt: row.CreatedAt,
				Title: row.Title,
				Price: int(row.Price),
				Amount: int(row.Amount),
				AlbumID: row.AlbumID,
				FileID: uuid.NullUUID{UUID: row.FileID.UUID, Valid: !row.FileID.UUID.IsNil()},
		}

		// add file
		if !row.PackFileID.UUID.IsNil() {
			pack.File = &File{
				ID: row.PackFileID,
				Name: row.PackFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(row.PackFilePath.String, ""),
			}
		}

		// add album
		if !row.AlbumAlbumID.IsNil() {
			pack.Album = &Album{
				ID: row.AlbumAlbumID,
				Title: row.AlbumTitle,
			}
		}

		// unmarshal PackRarities
		var packRarities []PackRarity
		json.Unmarshal(row.PackRarities, &packRarities)
		pack.Rarities = packRarities

		packs = append(packs, pack)
	}

	return ShopPacksResponse{
		Packs: packs,
	}
}

func castToShopBundlesResponse(rows []database.GetShopBundlesRow) ShopBundlesResponse {
	if rows == nil || len(rows) <= 0 {
		return ShopBundlesResponse{
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

	return ShopBundlesResponse{
		Bundles: bundles,
	}
}
