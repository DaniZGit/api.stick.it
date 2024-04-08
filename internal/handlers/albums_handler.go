package handlers

import (
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
		return ctx.ErrorResponse(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildAlbumResponse(albums, nil))
}

////////////////////////////
/* GET - "/albums/:id" */
///////////////////////////
func GetAlbum(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AlbumGetRequest)
	if err := c.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	album, err := ctx.Queries.GetAlbum(ctx.Request().Context(), uuid.FromStringOrNil(a.ID))
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildAlbumResponse(album, nil))
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
	file := database.File{}

	// get uploaded file
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "albums", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}
	
	// create album
	album, err := ctx.Queries.CreateAlbum(ctx.Request().Context(), database.CreateAlbumParams{
		ID: newUUID,
		Title: a.Title,
		DateFrom: pgtype.Timestamp{Time: utils.StringToTime(a.DateFrom, true), Valid: true},
		DateTo: pgtype.Timestamp{Time: utils.StringToTime(a.DateTo, false), Valid: true},
		Featured: pgtype.Bool{Bool: a.Featured, Valid: true},
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildAlbumResponse(album, &file))
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

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
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
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
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

	return ctx.JSON(http.StatusCreated, data.BuildAlbumResponse(album, &file))
}

///////////////////////////////
/* DELETE - "/albums/:id" */
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

////////////////////////////////
/* GET - "/albums/:id/packs" */
////////////////////////////////
func GetAlbumPacks(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AlbumPacksGetRequest)
	if err := ctx.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	packs, err := ctx.Queries.GetAlbumPacks(ctx.Request().Context(), a.AlbumID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPackResponse(packs, &database.File{}))
}