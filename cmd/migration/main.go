package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DaniZGit/api.stick.it/internal/environment"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

func main() {
	// loads .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file:", err)
	}
	
	dbDriver := environment.GooseDriver()
	dbString := environment.GooseDSN()
	
	args := os.Args
	if (len(args) < 2) {
		fmt.Println("Please specify a goose command: up | down | reset | etc...")
		return
	}

	// such as up | down | reset | etc...
	command := args[1]

	db, err := goose.OpenDBWithDriver(dbDriver, dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, "internal/db/schema/", arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}