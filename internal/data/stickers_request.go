package data

import (
	"github.com/gofrs/uuid"
)

type StickerCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	Type string `json:"type" form:"type" validate:"required,oneof=image audio gif"`
	Top float32 `json:"top" form:"top" validate:"numeric,omitempty"`
	Left float32 `json:"left" form:"left" validate:"numeric,omitempty"`
	Width float32 `json:"width" form:"width" validate:"numeric,omitempty"`
	Height float32 `json:"height" form:"height" validate:"numeric,omitempty"`
	Numerator int `json:"numerator" form:"numerator" validate:"numeric,omitempty"`
	Denominator int `json:"denominator" form:"denominator" validate:"numeric,omitempty"`
	Rotation float32 `json:"rotation" form:"rotation" validate:"numeric,omitempty"`
	PageID uuid.UUID `json:"page_id" form:"page_id" validate:"required"`
	RarityID uuid.UUID `json:"rarity_id" form:"rarity_id"`
	StickerID uuid.UUID `json:"sticker_id" form:"sticker_id"`
}

type PageStickersGetRequest struct {
	PageID uuid.UUID `param:"page_id" validate:"required"`
}

type StickerUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title" validate:"required"`
	Type string `json:"type" form:"type" validate:"required,oneof=image audio gif"`
	Top float32 `json:"top" form:"top" validate:"numeric,omitempty"`
	Left float32 `json:"left" form:"left" validate:"numeric,omitempty"`
	Width float32 `json:"width" form:"width" validate:"numeric,omitempty"`
	Height float32 `json:"height" form:"height" validate:"numeric,omitempty"`
	Numerator float32 `json:"numerator" form:"numerator" validate:"numeric,omitempty"`
	Denominator float32 `json:"denominator" form:"denominator" validate:"numeric,omitempty"`
	Rotation float32 `json:"rotation" form:"rotation" validate:"numeric,omitempty"`
	FileID uuid.UUID `json:"file_id" form:"file_id"`
}

type StickerRaritiesGetRequest struct {
	StickerID uuid.UUID `param:"id" validate:"required"`
}

type StickerDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}