package data

import "github.com/gofrs/uuid"

type BundleCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	Price float32 `json:"price" form:"price" validate:"required"`
	Tokens int `json:"tokens" form:"tokens" validate:"required"`
	Bonus *int `json:"bonus" form:"bonus" validate:"required"`
}

type BundleUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title" validate:"required"`
	Price float32 `json:"price" form:"price" validate:"required"`
	Tokens int `json:"tokens" form:"tokens" validate:"required"`
	Bonus *int `json:"bonus" form:"bonus" validate:"required"`
	FileID string `json:"file_id" form:"file_id"`
}

type BundleDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}