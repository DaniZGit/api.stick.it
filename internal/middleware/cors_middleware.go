package middleware

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/environment"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS() middleware.CORSConfig {
	url := environment.FrontendUrl()

	return middleware.CORSConfig{
		Skipper: middleware.DefaultSkipper,
		AllowOrigins: []string{url},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
}