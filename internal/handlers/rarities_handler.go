package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func CreateRarity(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	r := new(data.RarityCreateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(r); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// create role
	rarity, err := ctx.Queries.CreateRarity(ctx.Request().Context(), database.CreateRarityParams{
		ID: uuid.Must(uuid.NewV4()),
		Title: r.Title,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildRarityResponse(rarity, data.Metadata{}))
}

func GetRarities(c echo.Context) error {
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

	rarities, err := ctx.Queries.GetRarities(ctx.Request().Context(), database.GetRaritiesParams{
		Limit: int32(limit),
		Offset: int32(page),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// build metadata
	metadata := data.Metadata{}
	if len(rarities) > 0 {
		metadata.CurrPage = int32(page)
		metadata.PageSize = int32(limit)
		metadata.TotalRecords = int32(rarities[0].TotalRows)
		metadata.FirstPage = 1
		metadata.LastPage = int32(math.Ceil(float64(metadata.TotalRecords) / float64(metadata.PageSize)))
	}

	return ctx.JSON(http.StatusCreated, data.BuildRarityResponse(rarities, metadata))
}