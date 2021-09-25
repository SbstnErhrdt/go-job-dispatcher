package connections

import (
	"github.com/SbstnErhrdt/env"
	"github.com/SbstnErhrdt/go-gorm-all-sql/pkg/sql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var SQLClient *gorm.DB

// ConnectToSQL inits a sql client
func ConnectToSQL() {
	// check if the necessary sql variables are set
	env.CheckRequiredEnvironmentVariables(
		"SQL_TYPE",
		"SQL_HOST",
		"SQL_USER",
		"SQL_PASSWORD",
		"SQL_PORT",
		"SQL_DATABASE",
	)
	log.Println("try to connect to sql database")
	client, err := sql.ConnectToDatabase()
	if err != nil {
		log.Fatal("failed to connected to sql database")
		return
	}
	SQLClient = client
	log.Info("successfully connected to sql database")
	return
}

// CloseSQLConnection closes the connection to sql
func CloseSQLConnection() {
	log.Println("try to close connection to sql database")
	db, err := SQLClient.DB()
	if err != nil {
		log.Fatal("failed to get sql db")
		return
	}
	err = db.Close()
	if err != nil {
		log.Fatal("failed to close connection")
		return
	}
	log.Info("successfully closed connection")
}
