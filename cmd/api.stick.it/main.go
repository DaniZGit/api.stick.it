package main

import (
	"context"
	"log"

	"github.com/DaniZGit/api.stick.it/environment"
	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/db"
	api_middleware "github.com/DaniZGit/api.stick.it/internal/middleware"
	"github.com/DaniZGit/api.stick.it/internal/routes"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// loads .env file
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error while loading .env file:", err)
  }

	// create a new Echo instance
	e := echo.New()

	// create db instance
	conn, queries := db.Init()
	defer conn.Close(context.Background())

	// use extended context middleware
	e.Use(app.ExtendedContext(queries))

	// set CORS configuration
	e.Use(middleware.CORSWithConfig(api_middleware.CORS()))

	// set custom validator
	e.Validator = &app.ApiValidator{
		Validator: app.InitValidator(),
	}

	// add routes/endpoints
	routes.V1(e)

	// expose assets folder
	e.Static("/assets", "assets/public")

	// start Echo server
	startServer(e)
}

func startServer(e *echo.Echo) {
	serverUrl := environment.ServerUrl()

	// Start
	e.Logger.Fatal(e.Start(serverUrl))
}