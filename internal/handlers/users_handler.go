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
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

////////////////////////
/* GET - "/users/:id" */
////////////////////////
func GetUser(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	claims := auth.GetClaimsFromToken(*ctx)

	user, err := ctx.Queries.GetUserByID(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	token := ctx.Get("user").(*jwt.Token)

	return ctx.JSON(http.StatusOK, data.CastToUserByIDResponse(user, token.Raw))
}

func UpdateUser(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserUpdateRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	claims := auth.GetClaimsFromToken(*ctx)

	_, err := ctx.Queries.UpdateUser(ctx.Request().Context(), database.UpdateUserParams{
		ID: claims.UserID,
		Description: u.Description,
		AvatarID: uuid.NullUUID{UUID: u.AvatarID, Valid: !u.AvatarID.IsNil()},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	user, err := ctx.Queries.GetUserByID(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	token := ctx.Get("user").(*jwt.Token)

	return ctx.JSON(http.StatusOK, data.CastToUserByIDResponse(user, token.Raw))
}

func GetUserAlbums(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserAlbumsGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	albums, err := ctx.Queries.GetUserAlbums(ctx.Request().Context(), u.ID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusCreated, data.BuildAlbumResponse(albums, nil))
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

	userStickers, err := ctx.Queries.GetUserStickersForAlbum(ctx.Request().Context(), database.GetUserStickersForAlbumParams{
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
	userStickers := []data.UserSticker{}
	for _, sticker := range stickers {
		userSticker, err := qtx.CreateUserSticker(ctx.Request().Context(), database.CreateUserStickerParams{
			ID: uuid.Must(uuid.NewV4()),
			UserID: u.ID,
			StickerID: sticker.ID,
			Amount: 1,
			Sticked: false,
		})
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}

		userStickers = append(userStickers, data.UserSticker{
			ID: userSticker.ID,
			UserID: userSticker.UserID,
			StickerID: userSticker.StickerID,
			Amount: int(userSticker.Amount),
			Sticked: userSticker.Sticked,
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
		})
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

func StickUserSticker(c echo.Context) error  {
	ctx := c.(*app.ApiContext)

	u := new(data.StickUserStickerRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userSticker, err := ctx.Queries.StickUserSticker(ctx.Request().Context(), database.StickUserStickerParams{
		UserID: u.ID,
		StickerID: u.StickerID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	
	return ctx.JSON(http.StatusOK, data.BuildStickerResponse(userSticker, &database.File{}, &database.Rarity{}))
}

///////////////////////////////
/* POST - "/users/free-pack" */
///////////////////////////////
func ClaimUserFreePack(c echo.Context) error  {
	ctx := c.(*app.ApiContext)

	u := new(data.ClaimUserFreePackRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	claims := auth.GetClaimsFromToken(*ctx)
	// claim free pack - decremenet available_free_packs column by 1
	user, err := ctx.Queries.ClaimUserFreePack(ctx.Request().Context(), claims.UserID)
	if err != nil { // if no rows are returning, means the user didn't have any free packs available
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// check if user had all freebies, if they did we need to reset obtain date to current date
	if user.AvailableFreePacks == 2-1 {
		user, err = ctx.Queries.ResetUserFreePackDate(ctx.Request().Context(), claims.UserID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	// add selected pack to the user
	userPack, err := ctx.Queries.CreateUserPack(ctx.Request().Context(), database.CreateUserPackParams{
		ID: uuid.Must(uuid.NewV4()),
		UserID: claims.UserID,
		PackID: u.PackID,
		Amount: 1,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	
	return ctx.JSON(http.StatusOK, echo.Map{
		"user": data.CastToUserResponse(user, claims.ID).User,
		"user_pack": data.BuildPackResponse(userPack, &database.File{}).(data.UserPackResponse).UserPack,
	})
}

func GetUserAuctionStickers(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserAuctionStickersGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userStickers, err := ctx.Queries.GetUserAuctionStickers(ctx.Request().Context(), u.ID)

	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildStickerResponse(userStickers, &database.File{}, &database.Rarity{}))
}

/////////////////////////
/* /users/:id/progress */
/////////////////////////
func GetUserProgress(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	completedAlbumsCount := 0
	claims := auth.GetClaimsFromToken(*ctx)
	completedAlbums, err := ctx.Queries.GetUserCompletedAlbumsCount(ctx.Request().Context(), claims.UserID)
	if err != nil && err != pgx.ErrNoRows {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	if err != pgx.ErrNoRows {
		completedAlbumsCount = len(completedAlbums)
	}

	foundStickersCount := 0
	stickersCount, err := ctx.Queries.GetUserFoundStickersCount(ctx.Request().Context(), claims.UserID)
	if err != nil && err != pgx.ErrNoRows {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}
	if err != pgx.ErrNoRows {
		foundStickersCount = int(stickersCount)
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"progress": echo.Map{
			"completed_albums_count": completedAlbumsCount,
			"found_stickers_count": foundStickersCount,
		},
	})
}