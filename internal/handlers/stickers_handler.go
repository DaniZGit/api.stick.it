package handlers

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/utils"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

////////////////////////////
/* POST - "/stickers	  " */
////////////////////////////
func CreateSticker(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	s := new(data.StickerCreateRequest)
	if err := ctx.Bind(s); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(s); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	newUUID := uuid.Must(uuid.NewV4())

	// get uploaded file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "stickers", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	sticker, err := ctx.Queries.CreateSticker(ctx.Request().Context(), database.CreateStickerParams{
		ID: newUUID,
		Title: s.Title,
		Type: s.Type,
		Top: utils.FloatToPgNumeric(s.Top, 0),
		Left: utils.FloatToPgNumeric(s.Left, 0),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
		PageID: s.PageID,
		RarityID: s.RarityID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildStickerResponse(sticker, &file))
}

////////////////////////////
/* POST - "/stickers	  " */
////////////////////////////
func GetPageStickers(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	s := new(data.PageStickersGetRequest)
	if err := ctx.Bind(s); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	stickers, err := ctx.Queries.GetPageStickers(ctx.Request().Context(), s.PageID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildStickerResponse(stickers, &database.File{}))
}