package data

import (
	"github.com/gofrs/uuid"
)

type StickerCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	Type string `json:"type" form:"type" validate:"required,oneof=image audio gif"`
	Top float32 `json:"top" form:"top" validate:"numeric,omitempty"`
	Left float32 `json:"left" form:"left" validate:"numeric,omitempty"`
	PageID uuid.UUID `json:"page_id" form:"page_id" validate:"required"`
	RarityID uuid.UUID `json:"rarity_id" form:"rarity_id" validate:"required"`
}

type PageStickersGetRequest struct {
	PageID uuid.UUID `param:"page_id" validate:"required"`
}

type UpdateStickerRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title" validate:"required"`
	Type string `json:"type" form:"type" validate:"required,oneof=image audio gif"`
	Top float32 `json:"top" form:"top" validate:"numeric,omitempty"`
	Left float32 `json:"left" form:"left" validate:"numeric,omitempty"`
	RarityID uuid.UUID `json:"rarity_id" form:"rarity_id" validate:"required"`
	FileID string `json:"file_id" form:"file_id"`
}
