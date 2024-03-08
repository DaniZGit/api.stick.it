package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type FileResponse struct {
	ID pgtype.UUID `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
}

func CastToFileResponse(f database.File) FileResponse {
	return FileResponse{
		ID: f.ID,
		Name: f.Name,
	}
}