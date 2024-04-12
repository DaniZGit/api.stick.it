package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Pack struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title string `json:"title"`
	Price int `json:"price"`
	Amount int `json:"amount"`
	AlbumID uuid.UUID `json:"album_id"`
	FileID uuid.UUID `json:"file_id"`
	File *File `json:"file"`
	Rarities []PackRarity `json:"rarities"`
}

type PackRarity struct {
	ID        uuid.UUID        `json:"id"`
	PackID uuid.UUID `json:"pack_id"`
	RarityID uuid.NullUUID `json:"rarity_id"`
	DropChance pgtype.Numeric `json:"drop_chance"`
}

type PackResponse struct {
	Pack Pack `json:"pack"`
}

type PacksResponse struct {
	Packs []Pack `json:"packs"`
}

type PackRaritiesResponse struct {
	PackRarities []PackRarity `json:"pack_rarities"`
}

type PackRarityResponse struct {
	PackRarity PackRarity `json:"pack_rarity"`
}

func BuildPackResponse(packRows interface{}, file *database.File) any {
	switch value := packRows.(type) {
		case database.Pack:
			return PackResponse{
				Pack: Pack{
					ID: value.ID,
					CreatedAt: value.CreatedAt,
					Title: value.Title,
					Price: int(value.Price),
					Amount: int(value.Amount),
					AlbumID: value.AlbumID,
					File: &File{
						ID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
						Name: file.Name,
						Url: assetmanager.GetPublicAssetsFileUrl(file.Path, ""),
					},
				},
			}
		case []database.GetAlbumPacksRow:
			return castToAlbumPacksResponse(value)
		case database.PackRarity:
			return PackRarityResponse{
				PackRarity: PackRarity{
					ID: value.ID,
					PackID: value.PackID,
					RarityID: value.RarityID,
					DropChance: value.DropChance,
				},
			}
		case []database.PackRarity:
			return castToPackRaritiesResponse(value)
	}

	return AlbumResponse{}
}

func castToAlbumPacksResponse(packsRows []database.GetAlbumPacksRow) PacksResponse {
	if packsRows == nil || len(packsRows) <= 0 {
		return PacksResponse{
			Packs: []Pack{},
		}
	}

	packs := []Pack{}
	for _, packsRow := range packsRows {
		pack := Pack{
			ID: packsRow.ID,
			CreatedAt: packsRow.CreatedAt,
			Title: packsRow.Title,
			Price: int(packsRow.Price),
			Amount: int(packsRow.Amount),
			AlbumID: packsRow.AlbumID,
			FileID: packsRow.FileID.UUID,
		}

		// add file
		if !packsRow.PackFileID.UUID.IsNil() {
			pack.File = &File{
				ID: packsRow.PackFileID,
				Name: packsRow.PackFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(packsRow.PackFilePath.String, ""),
			}
		}

		packs = append(packs, pack)
	}
	
	return PacksResponse{
		Packs: packs,
	}
}

func castToPackRaritiesResponse(packRaritiesRows []database.PackRarity) PackRaritiesResponse {
	if packRaritiesRows == nil || len(packRaritiesRows) <= 0 {
		return  PackRaritiesResponse{
			PackRarities: []PackRarity{},
		}
	}

	packRarities := []PackRarity{}
	for _, packRaritiesRow := range packRaritiesRows {
		packRarity := PackRarity{
			ID: packRaritiesRow.ID,
			PackID: packRaritiesRow.PackID,
			RarityID: packRaritiesRow.RarityID,
			DropChance: packRaritiesRow.DropChance,
		}

		packRarities = append(packRarities, packRarity)
	}

	return PackRaritiesResponse{
		PackRarities: packRarities,
	}
}
