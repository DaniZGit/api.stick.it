package app

import (
	"errors"

	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

type ApiContext struct{
	echo.Context
	Queries *database.Queries
}

func ExtendedContext(q *database.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create an ApiContext with the database queries
			cc := &ApiContext{
				Context: c,
				Queries: q,
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
