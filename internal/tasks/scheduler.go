package tasks

import (
	"fmt"

	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/ws"
	"github.com/go-co-op/gocron/v2"
)

func InitScheduler(queries *database.Queries, hubs *ws.HubModels) {
	fmt.Println("Starting task scheduler...")
	go startScheduler(queries, hubs)
}

func startScheduler(queries *database.Queries, hubs *ws.HubModels) {
	// create a scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		fmt.Println("Scheduler init error", err)
		return
	}
	defer func() { _ = scheduler.Shutdown() }()

	// add jobs to the scheduler
	addTasks(scheduler, queries, hubs)

	// start the scheduler
	scheduler.Start()

	// block until you are ready to shut down - in my case, block forever
	select {}
}

func addTasks(scheduler gocron.Scheduler, queries *database.Queries, hubs *ws.HubModels) {
	_, err := scheduler.NewJob(freebiesTask(queries))
	if err != nil {
		// handle error
		fmt.Println("error while creating freebies task", err)
		return
	}

	_, err = scheduler.NewJob(markCompletedAuctionsTask(queries, hubs))
	if err != nil {
		// handle error
		fmt.Println("error while creating completed auctions task", err)
		return
	}

	// fmt.Printf("Started job with id %v\n", job.ID())
}