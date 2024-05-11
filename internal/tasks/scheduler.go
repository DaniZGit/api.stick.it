package tasks

import (
	"fmt"

	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/go-co-op/gocron/v2"
)

func InitScheduler(queries *database.Queries) {
	fmt.Println("Starting task scheduler...")
	go startScheduler(queries)
}

func startScheduler(queries *database.Queries) {
	// create a scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		fmt.Println("Scheduler init error", err)
		return
	}
	defer func() { _ = scheduler.Shutdown() }()

	// add jobs to the scheduler
	addTasks(scheduler, queries)

	// start the scheduler
	scheduler.Start()

	// block until you are ready to shut down - in my case, block forever
	select {}
}

func addTasks(scheduler gocron.Scheduler, queries *database.Queries) {
	_, err := scheduler.NewJob(freebiesTask(queries))
	if err != nil {
		// handle error
		fmt.Println("error while creating job", err)
		return
	}

	// fmt.Printf("Started job with id %v\n", job.ID())
}