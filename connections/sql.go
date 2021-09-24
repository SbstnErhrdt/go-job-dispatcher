package connections

import (
	"github.com/SbstnErhrdt/env"
	"github.com/SbstnErhrdt/go-gorm-all-sql/pkg/sql"
	"gorm.io/gorm"
	"log"
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
	log.Println("Try to connect to sql database")
	client, err := sql.ConnectToDatabase()
	if err != nil {
		log.Println("Failed to connected to sql database")
		panic(err)
		return
	}
	SQLClient = client
	log.Println("Successfully connected to sql database")
	return
}

func CloseSQLConnection() {
	log.Println("Try to close connection to sql database")
	db, err := SQLClient.DB()
	if err != nil {
		log.Println("Failed to get sql db")
		panic(err)
		return
	}
	err = db.Close()
	if err != nil {
		log.Println("Failed to close connection")
		panic(err)
		return
	}
	log.Println("Successfully closed connection")
}
