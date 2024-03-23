package db

import (
	"context"
	"fmt"
	"log"

	"github.com/DaniZGit/api.stick.it/environment"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
	Creates a new DB instance
*/
func Init() (*pgxpool.Pool, *database.Queries)  {
	fmt.Print("Connecting to db...")

	// Connects to DB
	conn := newPoolConnection()
	fmt.Print("Connected\n")

	// Returns pool connection and db queries model
	return conn, database.New(conn)
}

func newPoolConnection() *pgxpool.Pool {
	// Create db pool connection
	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig())
	if err!=nil {
	 log.Fatal("Error while creating connection to the database!")
	} 

	return connPool
}

/*
	Creates a new DB connection
*/
func dbConfig() *pgxpool.Config {
	dbCredentials := environment.DbCredentials()

	// Creates a new db config
	dbConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbCredentials.User,
			dbCredentials.Pass,
			dbCredentials.Host,
			dbCredentials.Port,
			dbCredentials.Name,
			dbCredentials.SSL,
		),
	)
	if err != nil {
		log.Fatal("Failed to create db config:", err)
	}

	// read values from .env
	dbConfigEnv := environment.DbConfig()

	// set db config values
	dbConfig.MaxConns = dbConfigEnv.MaxConns
	dbConfig.MinConns = dbConfigEnv.MinConns
	dbConfig.MaxConnLifetime = dbConfigEnv.MaxConnLifetime
	dbConfig.MaxConnIdleTime = dbConfigEnv.MaxConnIdleTime
	dbConfig.HealthCheckPeriod = dbConfigEnv.HealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = dbConfigEnv.ConnectTimeout

	return dbConfig
}