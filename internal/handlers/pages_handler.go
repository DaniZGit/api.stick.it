package handlers

import (
	"net/http"
	"os"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func CreatePage(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PageCreateRequest)
	if err := ctx.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	newUUID := uuid.Must(uuid.NewV4())

	// get uploaded file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "pages", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	page, err := ctx.Queries.CreatePage(ctx.Request().Context(), database.CreatePageParams{
		ID: newUUID,
		SortOrder: p.SortOrder,
		AlbumID: p.AlbumID,
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPageResponse(page, &file))
}

func GetPage(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PageGetRequest)
	if err := ctx.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	page, err := ctx.Queries.GetPage(ctx.Request().Context(), uuid.FromStringOrNil(p.ID))
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPageResponse(page, &database.File{}))
}

////////////////////////
/* PUT - "/pages/:id" */
////////////////////////
func UpdatePage(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// parse data from body
	p := new(data.PageUpdateRequest)
	if err := c.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(p); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}
	
	// check for file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		// new file
		fileUUID := uuid.Must(uuid.NewV4())
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "pages", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(p.FileID)
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
	}
	
	// update page
	page, err := ctx.Queries.UpdatePage(ctx.Request().Context(), database.UpdatePageParams{
		ID: p.ID,
		SortOrder: int32(p.SortOrder),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPageResponse(page, &file))
}

///////////////////////////
/* DELETE - "/pages/:id" */
///////////////////////////
func DeletePage(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PageDeleteRequest)
	if err := c.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(p); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	// delete page from DB
	page, err := ctx.Queries.DeletePage(ctx.Request().Context(), p.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// delete page file from disk
	_ = os.Remove(assetmanager.GetAssetsFileUrl(page.FileID.UUID.String(), ""))

	return ctx.JSON(http.StatusOK, data.BuildPageResponse(page, &database.File{}))
}