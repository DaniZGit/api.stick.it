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
	"github.com/gofrs/uuid"
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

	id := ctx.Param("id")
	album, err := ctx.Queries.GetAlbum(ctx.Request().Context(), uuid.FromStringOrNil(id))
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

	newUUID := uuid.Must(uuid.NewV4())
	fileUUID := uuid.UUID{}

	// get uploaded file
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err := assetmanager.CreateFileWithUUID(f, ctx, "albums", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}

		fileUUID = file.ID
	}
	
	// create album
	_, err = ctx.Queries.CreateAlbum(ctx.Request().Context(), database.CreateAlbumParams{
		ID: newUUID,
		Title: a.Title,
		DateFrom: pgtype.Timestamp{Time: utils.StringToTime(a.DateFrom, true), Valid: true},
		DateTo: pgtype.Timestamp{Time: utils.StringToTime(a.DateTo, false), Valid: true},
		Featured: pgtype.Bool{Bool: a.Featured, Valid: true},
		FileID: uuid.NullUUID{UUID: fileUUID, Valid: !fileUUID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// get album with files
	album, err := ctx.Queries.GetAlbum(ctx.Request().Context(), newUUID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAlbumResponse(album))
}

///////////////////////////////
/* PUT - "/albums/:id" */
///////////////////////////////
func UpdateAlbum(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// parse data from body
	a := new(data.AlbumUpdateRequest)
	if err := c.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	// check for file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		// new file
		fileUUID := uuid.Must(uuid.NewV4())
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "albums", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(a.FileID)
		file, err = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}
	
	// update album
	album, err := ctx.Queries.UpdateAlbum(ctx.Request().Context(), database.UpdateAlbumParams{
		ID: a.ID,
		Title: a.Title,
		DateFrom: utils.StringToPgTime(a.DateFrom, false),
		DateTo: utils.StringToPgTime(a.DateTo, false),
		Featured: pgtype.Bool{Bool: a.Featured, Valid: true},
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAlbumUpdateResponse(album, file))
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
	album, err := ctx.Queries.DeleteAlbum(ctx.Request().Context(), uuid.FromStringOrNil(a.ID))
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// delete album file from disk
	_ = os.Remove(assetmanager.GetAssetsFileUrl(album.FileID.UUID.String(), ""))

	return ctx.NoContent(http.StatusOK)
}