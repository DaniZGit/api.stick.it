package app

import (
	"errors"

	"github.com/DaniZGit/api.stick.it/environment"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/mailer"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ApiContext struct{
	echo.Context
	Queries *database.Queries
	DBPool *pgxpool.Pool
	Mailer mailer.Mailer
}

func ExtendedContext(dbPool *pgxpool.Pool, q *database.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			smtpConfig := environment.SMTPConfig()

			// Create an ApiContext with the database queries
			cc := &ApiContext{
				Context: c,
				Queries: q,
				DBPool: dbPool,
				Mailer: mailer.New(
					smtpConfig.Host, 
					smtpConfig.Port, 
					smtpConfig.Username, 
					smtpConfig.Password, 
					smtpConfig.Sender,
				),
			}
			
			// Call the next middleware or handler in the chain
			return next(cc)
		}
	}
}

func(ctx ApiContext) ErrorResponse(code int, err error) error {
	var pgErr *pgconn.PgError
	var validationErrs validator.ValidationErrors

	switch {
		// postgres db errors
		case errors.As(err, &pgErr):
			// map db error field to payload fields
			// db error example: duplicate key value violates unique constraint (username/email field)
			var mappedConstraintErrors = map[string]string {
				"users_username_key": "Username",
				"users_email_key": "Email",
				"albums_title_unique": "Title",
			}

			// return db error in the same format as validation errors 
			if field, ok := mappedConstraintErrors[pgErr.ConstraintName]; ok {
				return ctx.JSON(code, echo.Map{
					"type": "db",
					"errors": []ValidationError{
						{
							Field: field,
							Tag: 	 pgErr.Code,
						},
					},
				})
			}

		// payload validation errors
		case errors.As(err, &validationErrs):
			errors := []ValidationError{}
			for _, err := range validationErrs {
				e := ValidationError{
					Field: err.Field(),
					Tag:   err.Tag(),
				}

				errors = append(errors, e)
			}

			return ctx.JSON(code, echo.Map{
				"type": "validation",
				"errors": errors,
			})
	}

	// if none of the above were satisfied, return the 'unexpected' error
	return ctx.JSON(code, echo.Map{
		"type": "other",
		"error": err.Error(),
	})
}
