package routes

import (
	"github.com/DaniZGit/api.stick.it/internal/handlers"
	"github.com/DaniZGit/api.stick.it/internal/middleware"
	"github.com/labstack/echo/v4"
)

func V1(e *echo.Echo) {
	v1 := e.Group("/v1")

	v1.POST("/register", handlers.UserRegister)
	v1.POST("/login", handlers.UserLogin)

	// use JWT auth
	v1.Use(middleware.JwtAuth())
	v1.GET("/users/:id", handlers.GetUser)
	v1.GET("/albums", handlers.GetAlbums)
	v1.GET("/albums/:id", handlers.GetAlbum)
	v1.POST("/albums", handlers.CreateAlbum)
	v1.PUT("/albums/:id", handlers.UpdateAlbum)
	v1.DELETE("/albums/:id", handlers.DeleteAlbum)

	v1.POST("/pages", handlers.CreatePage)
}
