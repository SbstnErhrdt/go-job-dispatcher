package connections

import (
	"context"
	"github.com/SbstnErhrdt/env"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var RedisClient *redis.Client

// ConnectToRedis creates a new connection to redis
func ConnectToRedis() {
	// check if the necessary sql variables are set
	env.CheckRequiredEnvironmentVariables(
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_DATABASE",
	)
	// check if optional variables are present
	env.CheckOptionalEnvironmentVariables(
		"REDIS_PASSWORD",
	)
	log.Info("try to connect to redis")
	db, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		log.Println("failed to convert the provided redis db in the environment variable REDIS_DATABASE to an integer")
		log.Fatal(err)
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       db,                          // use default DB
	})
	RedisClient = client
	res := RedisClient.Ping(context.TODO())
	if res.Err() != nil {
		log.Fatal(res.Err())
		return
	}
	log.Info("successfully connected to redis database")
	return
}

// CloseRedisConnection closes a connection
func CloseRedisConnection() {
	err := RedisClient.Close()
	if err != nil {
		log.Error("failed to close redis connection")
		panic(err)
		return
	}
}
