package migrations

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"log"
)

func Run() {
	log.Println("Try to migrate")
	// Init the connections to the databases
	connections.ConnectToSQL()
	// Create the tables in the database
	err := connections.SQLClient.AutoMigrate(
		&job_dispatcher.Job{},
	)
	if err != nil {
		log.Print("Could not migrate", err)
	}
	log.Println("Migration successful")
}
