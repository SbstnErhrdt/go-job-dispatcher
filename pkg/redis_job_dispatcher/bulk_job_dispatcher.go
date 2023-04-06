package redis_job_dispatcher

import (
	"context"
	"encoding/json"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	log "github.com/sirupsen/logrus"
)

// Push adds a new job to the bulk job list
func Push(job job_dispatcher.Job) (err error) {
	jobJson, err := json.Marshal(job)
	if err != nil {
		log.WithError(err).Error("Could not marshal job")
		return err
	}
	resCDM := connections.RedisClient.LPush(context.TODO(), job.WorkerInstance, jobJson)
	if resCDM.Err() != nil {
		err = resCDM.Err()
		log.WithError(err).Error("Could not push job to redis")
		return
	}
	return
}

// Pop retrieves a job to the bulk job list
func Pop(instance string) (job job_dispatcher.Job, err error) {
	resCDM := connections.RedisClient.LPop(context.TODO(), instance)
	if resCDM.Err() != nil {
		err = resCDM.Err()
		log.WithError(err).Error("Could not pop job from redis")
		return
	}
	var jsonData []byte
	err = resCDM.Scan(&jsonData)
	if err != nil {
		log.WithError(err).Error("Could not scan job from redis")
		return
	}
	err = json.Unmarshal(jsonData, &job)
	if err != nil {
		log.WithError(err).Error("Could not unmarshal job from redis")
		return
	}
	return
}
