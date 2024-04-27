package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/utils"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

//////////////////////
/* GET - "/bundles" */
//////////////////////
func GetBundles(c echo.Context) error {
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

	bundles, err := ctx.Queries.GetBundles(ctx.Request().Context(), database.GetBundlesParams{
		Limit: int32(limit),
		Offset: int32(page),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// build metadata
	metadata := data.Metadata{}
	if len(bundles) > 0 {
		metadata.CurrPage = int32(page)
		metadata.PageSize = int32(limit)
		metadata.TotalRecords = int32(bundles[0].TotalRows)
		metadata.FirstPage = 1
		metadata.LastPage = int32(math.Ceil(float64(metadata.TotalRecords) / float64(metadata.PageSize)))
	}

	return ctx.JSON(http.StatusCreated, data.BuildBundlesResponse(bundles, metadata, &database.File{}))
}

///////////////////////
/* POST - "/bundles" */
///////////////////////
func CreateBundle(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	b := new(data.BundleCreateRequest)
	if err := ctx.Bind(b); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(b); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	newUUID := uuid.Must(uuid.NewV4())
	// get uploaded file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "bundles", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}
	
	bundle, err := ctx.Queries.CreateBundle(ctx.Request().Context(), database.CreateBundleParams{
		ID: newUUID,
		Title: b.Title,
		Price: utils.FloatToPgNumeric(b.Price, 0.0),
		Tokens: int32(b.Tokens),
		Bonus: int32(*b.Bonus),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildBundlesResponse(bundle, data.Metadata{}, &file))
}

//////////////////////////
/* PUT - "/bundles/:id" */
//////////////////////////
func UpdateBundle(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// parse data from body
	b := new(data.BundleUpdateRequest)
	if err := c.Bind(b); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(b); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}
	
	// check for file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		// new file
		fileUUID := uuid.Must(uuid.NewV4())
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "bundles", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(b.FileID)
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
	}
	
	// update pack
	bundle, err := ctx.Queries.UpdateBundle(ctx.Request().Context(), database.UpdateBundleParams{
		ID: b.ID,
		Title: b.Title,
		Price: utils.FloatToPgNumeric(b.Price, 0.0),
		Tokens: int32(b.Tokens),
		Bonus: int32(*b.Bonus),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildBundlesResponse(bundle, data.Metadata{}, &file))
}

/////////////////////////////
/* DELETE - "/bundles/:id" */
/////////////////////////////
func DeleteBundle(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	b := new(data.BundleDeleteRequest)
	if err := ctx.Bind(b); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	bundle, err := ctx.Queries.DeleteBundle(ctx.Request().Context(), b.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildBundlesResponse(bundle, data.Metadata{}, &database.File{}))
}