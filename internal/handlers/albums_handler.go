package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

/////////////////////
/* GET - "/albums" */
/////////////////////
func GetAlbums(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	l := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 12
	}

	albums, err := ctx.Queries.GetAlbums(ctx.Request().Context(), int32(limit))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, ctx.ErrorResponse(http.StatusNotFound, err))
	}

	return ctx.JSON(http.StatusCreated, data.CastToAlbumsResponse(albums))
}

////////////////////////////
/* GET - "/albums/:title" */
///////////////////////////
func GetAlbum(c echo.Context) error {
	ctx := c.(*app.ApiContext)
	title := ctx.Param("title")

	album, err := ctx.Queries.GetAlbum(ctx.Request().Context(), title)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, errors.New("album does not exist"))
	}

	return ctx.JSON(http.StatusOK, data.CastToAlbumResponse(album))
}

////////////////////////////
/* POST - "/albums" */
///////////////////////////
func CreateAlbum(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AlbumCreateRequest)
	if err := ctx.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	newUUID := uuid.New()
	fileUUID := pgtype.UUID{}

	// get uploaded file
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err := assetmanager.CreateFile(f, ctx, "albums", newUUID.String())
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}

		fileUUID = file.ID
	}
	
	// create album
	_, err = ctx.Queries.CreateAlbum(ctx.Request().Context(), database.CreateAlbumParams{
		ID: pgtype.UUID{Bytes: newUUID, Valid: true},
		Title: a.Title,
		DateFrom: pgtype.Timestamp{Time: utils.StringToTime(a.DateFrom, true), Valid: true},
		DateTo: pgtype.Timestamp{Time: utils.StringToTime(a.DateTo, false), Valid: true},
		FileID: fileUUID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// get album with files
	album, err := ctx.Queries.GetAlbum(ctx.Request().Context(), a.Title)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAlbumResponse(album))
}

///////////////////////////////
/* DELETE - "/albums/:title" */
///////////////////////////////
func DeleteAlbum(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AlbumDeleteRequest)
	if err := c.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	// delete album from DB
	album, err := ctx.Queries.DeleteAlbum(ctx.Request().Context(), a.Title)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// delete album file from disk
	_ = os.Remove(assetmanager.GetAssetsFileUrl(utils.UUIDToString(album.FileID), ""))

	return ctx.NoContent(http.StatusOK)
}