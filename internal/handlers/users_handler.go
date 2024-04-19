package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

////////////////////////
/* GET - "/users/:id" */
////////////////////////
func GetUser(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	claims := auth.GetClaimsFromToken(*ctx)

	return ctx.JSON(http.StatusOK, echo.Map{
		"user_id": claims.UserID,
		"role_id": claims.RoleID,
	})
}

func GetUserPacks(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserPacksGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userPacks, err := ctx.Queries.GetUserPacks(ctx.Request().Context(), database.GetUserPacksParams{
		UserID: u.ID,
		AlbumID: u.AlbumID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPackResponse(userPacks, &database.File{}))
}

func GetUserStickers(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserStickersGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userStickers, err := ctx.Queries.GetUserStickers(ctx.Request().Context(), database.GetUserStickersParams{
		UserID: u.ID,
		AlbumID: u.AlbumID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildStickerResponse(userStickers, &database.File{}, &database.Rarity{}))
}

func OpenUserPacks(c echo.Context) error  {
	ctx := c.(*app.ApiContext)

	u := new(data.UserPackOpenRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// start transaction
	tx, err := ctx.DBPool.Begin(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx.Request().Context())
	qtx := ctx.Queries.WithTx(tx)

	// pack opening logic
	pack, err := qtx.GetPackDetailed(ctx.Request().Context(), u.PackID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	var packRarities []database.PackRarity
	err = json.Unmarshal(pack.PackRarities, &packRarities)
	if err != nil || len(packRarities) <= 0 {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	numOfPacks := 1
	// check if user wants to open all packs
	if u.OpenAll {
		userPack, err := qtx.GetUserPack(ctx.Request().Context(), database.GetUserPackParams{
			UserID: u.ID,
			PackID: u.PackID,
		})
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}

		numOfPacks = int(userPack.Amount)
	}

	droppedStickers := make(map[uuid.NullUUID]int)
	dropChanceSum := 1.0
	// find the amount of stickers needed for each rarity
	for i := 0; i < int(pack.Amount) * numOfPacks; i++ {
		randomChance := rand.Float64() * dropChanceSum
		
		for _, packRarity := range packRarities {
			val, err := packRarity.DropChance.Float64Value()
			if err != nil {
				return ctx.ErrorResponse(http.StatusInternalServerError, err)
			}
			
			if randomChance <= val.Float64 {
				droppedStickers[packRarity.RarityID]++
				break;
			}

			randomChance -= val.Float64
		}
	}

	// get random stickers from db (based on rarity)
	stickers := []database.GetRandomStickersRow{}
	for val := range droppedStickers {
		// we want to get stickers for each pack one by one
		// so instead of doing 1 query to get all random stickers of that rarity, we do it one by one per pack
		// for example when album had 4 stickers and user opened 20 packs at once, the query had a LIMIT = 80, but that many stickers did not exist for that album,
		// so thats why i had to do them in their own groups --> this can be optimized further, using group by or json_agg to group stickers
		count := droppedStickers[val] / int(pack.Amount)
		if droppedStickers[val] % int(pack.Amount) > 0 { count++ }

		for i := 0; i < count; i++ {
			randomStickers, err := qtx.GetRandomStickers(ctx.Request().Context(), database.GetRandomStickersParams{
				AlbumID: u.AlbumID,
				RarityID: val,
				Limit: int32(min(droppedStickers[val], int(pack.Amount))),
			})
			if err != nil {
				return ctx.ErrorResponse(http.StatusInternalServerError, err)
			}

			stickers = append(stickers, randomStickers...)
		}
	}

	// add sticker to user in db
	mappedUserStickers := make(map[uuid.UUID]data.UserSticker)
	for _, sticker := range stickers {
		userSticker, err := qtx.CreateUserSticker(ctx.Request().Context(), database.CreateUserStickerParams{
			ID: uuid.Must(uuid.NewV4()),
			UserID: u.ID,
			StickerID: sticker.ID,
			Amount: 1,
		})
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}

		mappedUserStickers[sticker.ID] = data.UserSticker{
			ID: userSticker.ID,
			UserID: userSticker.UserID,
			StickerID: userSticker.StickerID,
			Amount: int(userSticker.Amount),
			Sticker: data.Sticker{
				ID: uuid.NullUUID{UUID: sticker.ID, Valid: !sticker.ID.IsNil()},
				CreatedAt: sticker.CreatedAt,
				Title: sticker.Title,
				Type: sticker.Type,
				Top: sticker.Top,
				Left: sticker.Left,
				Width: sticker.Width,
				Height: sticker.Height,
				Numerator: sticker.Numerator,
				Denominator: sticker.Denominator,
				Rotation: sticker.Rotation,
				PageID: sticker.PageID,
				RarityID: sticker.RarityID,
				FileID: sticker.FileID,
				File: &data.File{
					ID: uuid.NullUUID{UUID: sticker.StickerFileID.UUID, Valid: !sticker.StickerFileID.UUID.IsNil()},
					Name: sticker.StickerFileName.String,
					Url: assetmanager.GetPublicAssetsFileUrl(sticker.StickerFilePath.String, ""),
				},
				Rarity: &data.Rarity{
					ID: uuid.NullUUID{UUID: sticker.StickerRarityID.UUID, Valid: !sticker.StickerRarityID.UUID.IsNil()},
					Title: sticker.StickerRarityTitle.String,
				},
			},
		}
	}

	// decrement user pack's amount
	_, err = qtx.DecrementUserPack(ctx.Request().Context(), database.DecrementUserPackParams{
		UserID: u.ID,
		PackID: u.PackID,
		Amount: int32(numOfPacks),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// mapped values to slice
	userStickers := make([]data.UserSticker, 0, len(mappedUserStickers))
	for  _, value := range mappedUserStickers {
		userStickers = append(userStickers, value)
	}

	// randomize stickers
	rand.Shuffle(len(userStickers), func(i, j int) { userStickers[i], userStickers[j] = userStickers[j], userStickers[i] })

	// commit transaction
	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, data.UserStickersResponse{
		UserStickers: userStickers,
	})
}