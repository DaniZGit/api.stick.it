package data

import "github.com/gofrs/uuid"

type AvatarCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
}

type AvatarUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title"`
	FileID string `json:"file_id" form:"file_id"`
}

type AvatarDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}