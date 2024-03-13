package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

type File struct {
	ID uuid.NullUUID `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
}

type FileResponse struct {
	ID uuid.NullUUID `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
}

func CastToFileResponse(file database.File) FileResponse {
	return FileResponse{
		ID: uuid.NullUUID{UUID: file.ID, Valid: true},
		Name: file.Name,
	}
}