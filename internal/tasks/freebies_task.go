package tasks

import (
	"context"
	"fmt"
	"time"

	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/go-co-op/gocron/v2"
)

func freebiesTask(queries *database.Queries) (gocron.JobDefinition, gocron.Task) {
	cronDuration := gocron.DurationJob(
		1*time.Second,
	)
	cronTask := gocron.NewTask(
		func() {
			// do things
			err := queries.UpdateUsersFreePacks(context.Background(), 2)
			if err != nil {
				fmt.Println("Error while updating user free packs", err)
			}
		},
	)
	return cronDuration, cronTask
}
