package tests

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/api"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/scripts/migrations"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

var Server *httptest.Server
var WorkerUID uuid.UUID
var Client job_dispatcher.Client

func init() {
	// load env vars
	_ = godotenv.Load("../.env")
	// set up the db
	connections.ConnectToSQL()
	connections.ConnectToRedis()
	// run the migrations
	migrations.Run()
	// init the test server
	testServer := httptest.NewServer(api.InitServerEngine())
	Server = testServer
	// init uuid
	WorkerUID, _ = uuid.Parse("aaf38d15-0fb6-41b2-b46e-edddf25689ba")
	// init client
	c := job_dispatcher.NewClient(Server.URL, "test", WorkerUID)
	Client = c
}

func TestMain(m *testing.M) {
	// test the redis version
	if os.Getenv("TEST_ENV") == "REDIS" {
		// redis
		Client.UseBulk()
		log.Println("Start package redis test tests")
		_ = m.Run()
		log.Println("Done package redis test tests")
	} else {
		// test the sql version
		log.Println("Start package sql test tests")
		if !connections.SQLClient.Migrator().HasTable(&job_dispatcher.Job{}) {
			_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
		}
		_ = m.Run()
		_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
		log.Println("Done package sql test tests")
	}
}
