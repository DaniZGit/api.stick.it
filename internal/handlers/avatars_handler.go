package handlers

import (
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

///////////////////////
/* POST - "/avatars" */
///////////////////////
func CreateAvatar(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AvatarCreateRequest)
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
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "avatars", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}
	
	// create avatar
	avatar, err := ctx.Queries.CreateAvatar(ctx.Request().Context(), database.CreateAvatarParams{
		ID: newUUID,
		Title: a.Title,
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAvatarResponse(avatar, file))
}

//////////////////////
/* GET - "/avatars" */
//////////////////////
func GetAvatars(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	l := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 12
	}

	p := ctx.QueryParam("page")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}

	avatars, err := ctx.Queries.GetAvatars(ctx.Request().Context(), database.GetAvatarsParams{
		Limit: int32(limit),
		Offset: int32(limit * page),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// build metadata
	metadata := data.Metadata{}
	if len(avatars) > 0 {
		metadata.CurrPage = int32(page)
		metadata.PageSize = int32(limit)
		metadata.TotalRecords = int32(avatars[0].TotalRows)
		metadata.FirstPage = 1
		metadata.LastPage = int32(math.Ceil(float64(metadata.TotalRecords) / float64(metadata.PageSize)))
	}

	return ctx.JSON(http.StatusCreated, data.CastToAvatarsResponse(avatars, metadata))
}

//////////////////////////
/* PUT - "/avatars/:id" */
//////////////////////////
func UpdateAvatar(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// parse data from body
	a := new(data.AvatarUpdateRequest)
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
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "avatars", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(a.FileID)
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
	}
	
	// update avatar
	avatar, err := ctx.Queries.UpdateAvatar(ctx.Request().Context(), database.UpdateAvatarParams{
		ID: a.ID,
		Title: a.Title,
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAvatarResponse(avatar, file))
}

/////////////////////////////
/* DELETE - "/avatars/:id" */
/////////////////////////////
func DeleteAvatar(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AvatarDeleteRequest)
	if err := c.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	// delete avatar from DB
	avatar, err := ctx.Queries.DeleteAvatar(ctx.Request().Context(), a.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// delete avatar file from disk
	_ = os.Remove(assetmanager.GetAssetsFileUrl(avatar.FileID.UUID.String(), ""))

	return ctx.NoContent(http.StatusOK)
}