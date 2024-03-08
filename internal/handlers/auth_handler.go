package handlers

import (
	"errors"
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

////////////////////////
/* POST - "/register" */
////////////////////////
func UserRegister(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// bind payload to struct
	u := new(data.UserRegisterParams)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// validate payload
	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// use bcrypt to hash user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// create user
	user, err := ctx.Queries.CreateUser(ctx.Request().Context(), database.CreateUserParams{
		ID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username: u.Username,
		Email: u.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// generate a new jwt token and set cookie
	t, err := auth.CreateJwtToken(user)
	if err != nil {
		return ctx.ErrorResponse(http.StatusBadRequest, err)
	}

	// return user with token
	return ctx.JSON(
		http.StatusCreated,
		data.UserResponse{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			Username: user.Username,
			Email: user.Email,
			Token: t,
		},
	)
}

/////////////////////
/* POST - "/login" */
/////////////////////
func UserLogin(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	// strip payload body
	u := new(data.UserLoginParams)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// validate payload
	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// get user from DB
	user, err := ctx.Queries.GetUser(ctx.Request().Context(), u.Email)
	if err != nil {
		return ctx.ErrorResponse(http.StatusUnauthorized, errors.New("user with the provided email does not exist"))
	}

	// compare user's password with their hashed variant in DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return ctx.ErrorResponse(http.StatusUnauthorized, errors.New("email or password does not match"))
	}

	// generate a new jwt token and set cookie
	t, err := auth.CreateJwtToken(user)
	if err != nil {
		return ctx.ErrorResponse(http.StatusBadRequest, err)
	}

	// return user with token
	return ctx.JSON(
		http.StatusCreated,
		data.UserResponse{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			Username: user.Username,
			Email: user.Email,
			Token: t,
		},
	)
}
