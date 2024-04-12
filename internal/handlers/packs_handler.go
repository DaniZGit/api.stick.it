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

//////////////////////
/* POST - "/packs"  */
//////////////////////
func CreatePack(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PackCreateRequest)
	if err := ctx.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	newUUID := uuid.Must(uuid.NewV4())

	// get uploaded file
	file := database.File{}
	f, err := ctx.FormFile("file")
	if err == nil {
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "packs", newUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	pack, err := ctx.Queries.CreatePack(ctx.Request().Context(), database.CreatePackParams{
		ID: newUUID,
		AlbumID: p.AlbumID,
		Title: p.Title,
		Price: int32(p.Price),
		Amount: int32(p.Amount),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(pack, &file))
}

////////////////////////
/* PUT - "/packs/:id" */
////////////////////////
func UpdatePack(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// parse data from body
	p := new(data.PackUpdateRequest)
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
		file, err = assetmanager.CreateFileWithUUID(f, ctx, "packs", fileUUID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		// get current file, if any
		fileUUID := uuid.FromStringOrNil(p.FileID)
		file, _ = ctx.Queries.GetFile(ctx.Request().Context(), fileUUID)
	}
	
	// update pack
	pack, err := ctx.Queries.UpdatePack(ctx.Request().Context(), database.UpdatePackParams{
		ID: p.ID,
		Title: p.Title,
		Price: int32(p.Price),
		Amount: int32(p.Amount),
		FileID: uuid.NullUUID{UUID: file.ID, Valid: !file.ID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(pack, &file))
}

///////////////////////////
/* DELETE - "/packs/:id" */
///////////////////////////
func DeletePack(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PackDeleteRequest)
	if err := ctx.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	pack, err := ctx.Queries.DeletePack(ctx.Request().Context(), p.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(pack, &database.File{}))
}

/////////////////////////////////
/* GET - "/packs/:id/rarities" */
/////////////////////////////////
func GetPackRarities(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	p := new(data.PackRaritiesGetRequest)
	if err := ctx.Bind(p); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	packRarities, err := ctx.Queries.GetPackRarities(ctx.Request().Context(), p.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPackResponse(packRarities, &database.File{}))
}


/////////////////////////////
/* POST - "/pack-rarities" */
/////////////////////////////
func CreatePackRarity(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	pr := new(data.PackRarityCreateRequest)
	if err := ctx.Bind(pr); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	newUUID := uuid.Must(uuid.NewV4())
	pack, err := ctx.Queries.CreatePackRarity(ctx.Request().Context(), database.CreatePackRarityParams{
		ID: newUUID,
		PackID: pr.PackID,
		RarityID: uuid.NullUUID{UUID: pr.RarityID, Valid: !pr.RarityID.IsNil()},
		DropChance: utils.FloatToPgNumeric(pr.DropChance, 0.0),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(pack, &database.File{}))
}

////////////////////////////////
/* PUT - "/pack-rarities/:id" */
////////////////////////////////
func UpdatePackRarity(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	pr := new(data.PackRarityUpdateRequest)
	if err := ctx.Bind(pr); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	packRarity, err := ctx.Queries.UpdatePackRarity(ctx.Request().Context(), database.UpdatePackRarityParams{
		ID: pr.ID,
		DropChance: utils.FloatToPgNumeric(pr.DropChance, 0.0),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(packRarity, &database.File{}))
}

///////////////////////////////////
/* DELETE - "/pack-rarities/:id" */
///////////////////////////////////
func DeletePackRarity(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	pr := new(data.PackRarityDeleteRequest)
	if err := ctx.Bind(pr); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	
	packRarity, err := ctx.Queries.DeletePackRarity(ctx.Request().Context(), pr.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildPackResponse(packRarity, &database.File{}))
}