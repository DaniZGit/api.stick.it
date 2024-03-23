package seed

import (
	"context"
	"fmt"

	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/db"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

func main() {
	fmt.Print("Seeding db...");
	// create db instance
	dbPool, queries := db.Init()
	defer dbPool.Close()

	// seeders
	SeedRoles(queries)
	SeedUsers(queries)
	
	fmt.Print(" Complete\n")
}

func SeedRoles(queries *database.Queries) {
	// init 'admin' 'user' roles
	fmt.Print("Creating 'Admin' role...")
	_, err := queries.CreateRole(context.Background(), database.CreateRoleParams{
		ID: uuid.Must(uuid.NewV4()),
		Title: "Admin",
	})
	if err != nil {
		fmt.Print("Already exists\n")
	} else {
		fmt.Print("Created\n")
	}
}

func SeedUsers(queries *database.Queries) {
	fmt.Print("Creating admin user...")

	users, err := queries.GetUsers(context.Background(), database.GetUsersParams{
		Limit: 10,
		Offset: 0,
	})
	if err != nil {
		fmt.Println("Error fetching users")
	}

	if ( len(users) <= 0 ) {
		hash, err := auth.GeneratePassword("pass")
		if err != nil {
			fmt.Print("Error - pass\n")
		} else {
			role, err := queries.GetRoleByName(context.Background(), "Admin")
			if err != nil {
				fmt.Print("'Admin' role does not exist\n")
			}

			_, err = queries.CreateUser(context.Background(), database.CreateUserParams{
				ID: uuid.Must(uuid.NewV4()),
				Username: "admin",
				Email: "admin@gmail.com",
				Password: string(hash),
				RoleID: uuid.NullUUID{UUID: role.ID, Valid: !role.ID.IsNil()},
			})

			if err != nil {
				fmt.Print("Already exists\n")
			} else {
				fmt.Print("Created\n")
			}
		}
	} else {
		fmt.Println("Already exists")
	}
}