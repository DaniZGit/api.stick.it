package data

import "github.com/gofrs/uuid"

type RoleCreateRequest struct {
	RoleID uuid.UUID `json:"role_id" form:"role_id" validate:"required"`
	Title string `json:"title" form:"title" validate:"required"`
}

type RolesGetRequest struct {
	Limit int `json:"limit"`
	Page int `json:"page"`
}

type RoleGetByNameRequest struct {
	Title string `json:"title"`
}