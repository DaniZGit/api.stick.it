package routes

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/handlers"
	"github.com/DaniZGit/api.stick.it/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Global(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {return c.JSON(http.StatusOK, "pong")})
}

func V1(e *echo.Echo) {
	v1 := e.Group("/v1")

	v1.POST("/register", handlers.UserRegister)
	v1.POST("/login", handlers.UserLogin)
	v1.PUT("/confirmation", handlers.UserMailConfirmation)

	// use JWT auth
	v1.Use(middleware.JwtAuth())
	v1.GET("/users/:id", handlers.GetUser)
	v1.GET("/users/:id/packs", handlers.GetUserPacks)
	v1.GET("/users/:id/stickers", handlers.GetUserStickers)
	v1.POST("/users/:id/open-packs", handlers.OpenUserPacks)
	v1.PATCH("/users/:id/stick-sticker", handlers.StickUserSticker)
	v1.POST("/users/:id/free-pack", handlers.ClaimUserFreePack)

	v1.POST("/roles", handlers.CreateRole)
	v1.GET("/roles", handlers.GetRoles)
	v1.GET("/roles/:title", handlers.GetRoleByName)

	v1.GET("/albums", handlers.GetAlbums)
	v1.GET("/albums/:id", handlers.GetAlbum)
	v1.POST("/albums", handlers.CreateAlbum)
	v1.PUT("/albums/:id", handlers.UpdateAlbum)
	v1.DELETE("/albums/:id", handlers.DeleteAlbum)
	v1.GET("/albums/:id/packs", handlers.GetAlbumPacks)
	v1.GET("/albums/:id/pages", handlers.GetAlbumPages)

	v1.POST("/pages", handlers.CreatePage)
	v1.GET("/pages/:id", handlers.GetPage)
	v1.PUT("/pages/:id", handlers.UpdatePage)
	v1.DELETE("/pages/:id", handlers.DeletePage)
	// v1.GET("/pages/:page_id/stickers", handlers.GetPageStickers)

	v1.POST("/stickers", handlers.CreateSticker)
	v1.PUT("/stickers/:id", handlers.UpdateSticker)
	v1.DELETE("/stickers/:id", handlers.DeleteSticker)
	v1.GET("/stickers/:id/rarities", handlers.GetStickerRarities)

	v1.GET("/rarities", handlers.GetRarities)
	v1.POST("/rarities", handlers.CreateRarity)

	v1.POST("/packs", handlers.CreatePack)
	v1.PUT("/packs/:id", handlers.UpdatePack)
	v1.DELETE("/packs/:id", handlers.DeletePack)
	v1.GET("/packs/:id/rarities", handlers.GetPackRarities)
	
	v1.POST("/pack-rarities", handlers.CreatePackRarity)
	v1.PUT("/pack-rarities/:id", handlers.UpdatePackRarity)
	v1.DELETE("/pack-rarities/:id", handlers.DeletePackRarity)

	v1.GET("/bundles", handlers.GetBundles)
	v1.POST("/bundles", handlers.CreateBundle)
	v1.PUT("/bundles/:id", handlers.UpdateBundle)
	v1.DELETE("/bundles/:id", handlers.DeleteBundle)

	v1.GET("/shop/packs", handlers.GetShopPacks)
	v1.GET("/shop/bundles", handlers.GetShopBundles)

	// transactions
	v1.GET("/transactions/config", handlers.GetStripeConfig)
	v1.POST("/transactions/create-payment-intent", handlers.CreatePaymentIntent)
	v1.POST("/transactions/pack", handlers.BuyPack)
	v1.POST("/transactions/bundle", handlers.BuyBundle)

	// auction
	v1.POST("/auction/offers", handlers.CreateAuctionOffer)
	v1.GET("/auction/offers", handlers.GetAuctionOffers)
	v1.GET("/auction/bids", handlers.GetAuctionBids)
	v1.POST("/auction/bids", handlers.CreateAuctionBid)
}
