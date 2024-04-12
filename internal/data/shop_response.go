package data

import (
	"encoding/json"

	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
)

type ShopPacksResponse struct {
	Packs []Pack `json:"packs"`
}

func BuildShopResponse(rows interface{}) any {
	switch value := rows.(type) {
		case []database.GetShopPacksRow:
			return castToShopPacksResponse(value)
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
				FileID: row.FileID.UUID,
		}

		// add file
		if !row.PackFileID.UUID.IsNil() {
			pack.File = &File{
				ID: row.PackFileID,
				Name: row.PackFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(row.PackFilePath.String, ""),
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
