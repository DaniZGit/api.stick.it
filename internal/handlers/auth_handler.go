package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
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
	hashedPassword, err := auth.GeneratePassword(u.Password)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	newUUID := uuid.Must(uuid.NewV4())
	// create confirmation token
	confirmationToken := auth.GenerateConfirmationToken(newUUID)

	// create user
	user, err := ctx.Queries.CreateUser(ctx.Request().Context(), database.CreateUserParams{
		ID: newUUID,
		Username: u.Username,
		Email: u.Email,
		Password: string(hashedPassword),
		ConfirmationToken: pgtype.Text{String: confirmationToken, Valid: true},
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// send confirmation email
	go func() {
		err = ctx.Mailer.Send(user.Email, "account_confirmation_mail.tmpl", user)
		if err != nil {
			fmt.Println("Error while sending confirmation email", err)
		}
	}()

	// return user with token
	return ctx.NoContent(http.StatusCreated)
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
	err = auth.ValidatePassword(u.Password, user.Password)
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
		data.CastToUserResponse(user, t),
	)
}
