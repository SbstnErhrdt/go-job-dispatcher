package migrations

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	log "github.com/sirupsen/logrus"
)

func Run() {
	log.Info("Try to migrate")
	// Init the connections to the databases
	connections.ConnectToSQL()
	// Create the tables in the database
	err := connections.SQLClient.AutoMigrate(
		&job_dispatcher.Job{},
	)
	if err != nil {
		log.WithError(err).Error("Could not migrate the database")
		return
	}
	log.Info("Migration successful")
}
