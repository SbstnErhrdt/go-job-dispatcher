package main

import (
	"github.com/SbstnErhrdt/env"
	"github.com/SbstnErhrdt/go-job-dispatcher/api"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/scripts/migrations"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Try to start job dispatcher")
	// lod env file if there is one
	env.LoadEnvFiles()
	// connect to sql
	connections.ConnectToSQL()
	connections.ConnectToRedis()
	// migrate to the latest data structure
	migrations.Run()
	// start the background cleaning job
	go api.InitBackgroundJobs()
	// init the api server
	api.InitServer()
	log.Info("Job dispatcher started")
}
