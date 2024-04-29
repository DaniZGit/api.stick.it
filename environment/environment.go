package environment

import (
	"log"
	"os"
	"strconv"
	"time"
)

func AssetsUrl() string {
	val := readEnv("ASSETS_URL")
	return val
}

func ServerUrl() string {
	val := readEnv("SERVER_URL")
	return val
}

func FrontendUrl() string {
	val := readEnv("FRONTEND_URL")
	return val
}

func JwtSecret() string {
	val := readEnv("JWT_SECRET")
	return val
}

func StripeSecret() string {
	val := readEnv("STRIPE_SECRET_KEY")
	return val
}

func StripePublishableKey() string {
	val := readEnv("STRIPE_PUBLISHABLE_KEY")
	return val
}

type DbCredentialsType struct {
	Host string;
	Port string;
	User string;
	Pass string;
	Name string;
	SSL string;
}
func DbCredentials() DbCredentialsType {
	dbHost := readEnv("DB_HOST")
	dbPort := readEnv("DB_PORT")
	dbUser := readEnv("DB_USER")
	dbPass := readEnv("DB_PASS")
	dbName := readEnv("DB_NAME")

	dbSSL := readEnv("DB_SSL")
	if dbSSL != "enable" && dbSSL != "disable" {
		log.Fatal("DB_SSL must be either 'enable' or 'disable'")
	}

	return DbCredentialsType{
		Host: dbHost,
		Port: dbPort,
		User: dbUser,
		Pass: dbPass,
		Name: dbName,
		SSL: dbSSL,
	}
}

type DBConfigType struct {
	MaxConns int32
	MinConns int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout time.Duration
}
func DbConfig() DBConfigType {
	maxConnVar := readEnv("DB_MAX_CONN")
	maxConn, err := strconv.Atoi(maxConnVar)
	if err != nil || maxConn <= 0 {
		log.Fatal("DB_MAX_CONN value must be a positive number greater than 0")
	}

	minConnVar := readEnv("DB_MIN_CONN")
	minConn, err := strconv.Atoi(minConnVar)
	if err != nil || minConn < 0 {
		log.Fatal("DB_MIN_CONN value must be a positive number")
	}

	lifeTimeVar := readEnv("DB_LIFE_TIME")
	lifeTime, err := strconv.Atoi(lifeTimeVar)
	if err != nil || lifeTime < 0 {
		log.Fatal("DB_LIFE_TIME value must be a positive number (in seconds)")
	}

	idleTimeVar := readEnv("DB_IDLE_TIME")
	idleTime, err := strconv.Atoi(idleTimeVar)
	if err != nil || idleTime < 0 {
		log.Fatal("DB_IDLE_TIME value must be a positive number (in seconds)")
	}

	healthCheckVar := readEnv("DB_HC_PERIOD")
	healthCheck, err := strconv.Atoi(healthCheckVar)
	if err != nil || healthCheck < 0 {
		log.Fatal("DB_HC_PERIOD value must be a positive number greater than 0 (in seconds)")
	}

	timeoutVar := readEnv("DB_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutVar)
	if err != nil || timeout <= 0 {
		log.Fatal("DB_TIMEOUT value must be a positive number greater than 0 (in seconds)")
	}

	return DBConfigType{
		MaxConns: int32(maxConn),
		MinConns: int32(minConn),
		MaxConnLifetime: time.Second * time.Duration(maxConn),
		MaxConnIdleTime: time.Second * time.Duration(idleTime),
		HealthCheckPeriod: time.Second * time.Duration(healthCheck),
		ConnectTimeout: time.Second * time.Duration(timeout),
	}
}

func readEnv(name string) string {
	variable := os.Getenv(name)
	if variable == "" {
		log.Fatal(name, " is not found in the environment.")
	}

	return variable
}