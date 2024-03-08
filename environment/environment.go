package environment

import (
	"log"
	"os"
)

func AssetsUrl() string {
	val := os.Getenv("ASSETS_URL")
	if val == "" {
		log.Fatal("ASSETS_URL is not found in the environment.")
	}
	
	return val
}

func ServerUrl() string {
	val := os.Getenv("SERVER_URL")
	if val == "" {
		log.Fatal("SERVER_URL is not found in the environment.")
	}

	return val
}

func FrontendUrl() string {
	val := os.Getenv("FRONTEND_URL")
	if val == "" {
		log.Fatal("FRONTEND_URL is not found in the environment.")
	}

	return val
}

func JwtSecret() string {
	val := os.Getenv("JWT_SECRET")
	if val == "" {
		log.Fatal("JWT_SECRET is not found in the environment.")
	}

	return val
}

type DbCredentialsType struct {
	Host string;
	Port string;
	User string;
	Pass string;
	Name string;
}
func DbCredentials() DbCredentialsType {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST is not found in the environment.")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT is not found in the environment.")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER is not found in the environment.")
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		log.Fatal("DB_PASS is not found in the environment.")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not found in the environment.")
	}

	return DbCredentialsType{
		Host: dbHost,
		Port: dbPort,
		User: dbUser,
		Pass: dbPass,
		Name: dbName,
	}
}

func GooseDriver() string {
	dbDriver := os.Getenv("GOOSE_DRIVER")
	if dbDriver == "" {
		log.Fatal("GOOSE_DRIVER is not found in the environment.")
	}

	return dbDriver
}

func GooseDSN() string {
	dbstring := os.Getenv("GOOSE_DBSTRING")
	if dbstring == "" {
		log.Fatal("GOOSE_DBSTRING is not found in the environment.")
	}

	return dbstring
}
