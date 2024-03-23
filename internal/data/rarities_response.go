package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

type Rarity struct {
	ID     uuid.NullUUID        `json:"id"`
	Title  string           `json:"title"`
}

type RarityResponse struct {
	Rarity Rarity `json:"rarity"`
}

type RaritiesResponse struct {
	Metadata Metadata `json:"metadata"`
	Rarities []Rarity `json:"rarities"`
}

func BuildRarityResponse(rarityRows interface{}, metadata Metadata) any {
	switch value := rarityRows.(type) {
		case database.Rarity:
			return RarityResponse{
				Rarity: Rarity{
					ID: uuid.NullUUID{UUID: value.ID, Valid: !value.ID.IsNil()},
					Title: value.Title,
				},
			}
		case []database.GetRaritiesRow:
			return castToRaritiesResponse(value, metadata)
	}

	return RarityResponse{}
}

func castToRaritiesResponse(rarityRows []database.GetRaritiesRow, metadata Metadata) RaritiesResponse {
	if rarityRows == nil || len(rarityRows) <= 0 {
		return  RaritiesResponse{
			Rarities: []Rarity{},
		}
	}

	rarities := []Rarity{}
	for _, rarityRow := range rarityRows {
		rarity := Rarity{
			ID: uuid.NullUUID{UUID: rarityRow.ID, Valid: !rarityRow.ID.IsNil()},
			Title: rarityRow.Title,
		}

		rarities = append(rarities, rarity)
	}

	return RaritiesResponse{
		Metadata: metadata,
		Rarities: rarities,
	}
}