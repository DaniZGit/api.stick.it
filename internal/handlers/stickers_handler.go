package handlers

import (
	"net/http"
	"os"

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
		Width: utils.FloatToPgNumeric(s.Width, 0),
		Height: utils.FloatToPgNumeric(s.Height, 0),
		Numerator: int32(s.Numerator),
		Denominator: int32(s.Denominator),
		Rotation: utils.FloatToPgNumeric(s.Rotation, 0),
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

////////////////////////////
/* PUT - "/stickers/:id"  */
////////////////////////////
func UpdateSticker(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	s := new(data.UpdateStickerRequest)
	if err := ctx.Bind(s); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(s); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	// check for file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		// new file
		fileUUID := uuid.Must(uuid.NewV4())
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "stickers", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(s.FileID)
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
	}

	sticker, err := ctx.Queries.UpdateSticker(ctx.Request().Context(), database.UpdateStickerParams{
		ID: s.ID,
		Title: s.Title,
		Type: s.Type,
		Top: utils.FloatToPgNumeric(s.Top, 0),
		Left: utils.FloatToPgNumeric(s.Left, 0),
		Width: utils.FloatToPgNumeric(s.Width, 0),
		Height: utils.FloatToPgNumeric(s.Height, 0),
		Numerator: int32(s.Numerator),
		Denominator: int32(s.Denominator),
		Rotation: utils.FloatToPgNumeric(s.Rotation, 0),
		RarityID: s.RarityID,
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildStickerResponse(sticker, &file))
}

///////////////////////////////
/* DELETE - "/stickers/:id"  */
///////////////////////////////
func DeleteSticker(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	s := new(data.StickerDeleteRequest)
	if err := c.Bind(s); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// delete album from DB
	sticker, err := ctx.Queries.DeleteSticker(ctx.Request().Context(), s.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// delete album file from disk
	_ = os.Remove(assetmanager.GetAssetsFileUrl(sticker.FileID.UUID.String(), ""))

	return ctx.NoContent(http.StatusOK)
}