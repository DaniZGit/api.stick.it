package db

import (
	"context"
	"fmt"
	"log"

	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/environment"
	"github.com/jackc/pgx/v5"
)

/*
	Creates a new DB instance
*/
func Init() (*pgx.Conn, *database.Queries)  {
	// Connects to DB
	conn := newConnection()

	// Returns db queries model
	return conn, database.New(conn)
}

/*
	Creates a new DB connection
*/
func newConnection() *pgx.Conn {
	dbCredentials := environment.DbCredentials()

	environment.DbCredentials()

	// Creates a new connection to postgres database
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbCredentials.User,
			dbCredentials.Pass,
			dbCredentials.Host,
			dbCredentials.Port,
			dbCredentials.Name,
		),
	)

	if err != nil {
		log.Fatal("Error while connecting to db:", err)
	}

	return conn
}