package data

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Sticker struct {
	ID        uuid.NullUUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Title  string           `json:"title"`
	File *File `json:"file"`
}