package data

import "github.com/gofrs/uuid"

type PageCreateRequest struct {
	AlbumID uuid.UUID `json:"album_id" form:"album_id" validate:"required"`
	SortOrder int32 `json:"sort_order" form:"sort_order"`
}