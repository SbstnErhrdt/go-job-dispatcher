package main

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/redis_job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/sql_job_dispatcher"
	log "github.com/sirupsen/logrus"
)

func init() {
	connections.ConnectToSQL()
	connections.ConnectToRedis()
}

func main() {
	log.Info("started")
	cleanSql()
	cleanRedis()
	log.Info("done goodbye")
	return
}

func cleanSql() {
	log.Info("start sql clean")
	s := sql_job_dispatcher.SqlService{}
	err := s.Clean()
	if err != nil {
		log.Fatal("can not run procedure: err: ", err)
	}
	connections.CloseSQLConnection()
	log.Info("sql clean done")
}

func cleanRedis() {
	log.Info("start redis clean")
	s := redis_job_dispatcher.RedisService{}
	err := s.Clean()
	if err != nil {
		log.Fatal("can not run procedure: err: ", err)
	}
	connections.CloseRedisConnection()
	log.Info("redis clean done")
}
