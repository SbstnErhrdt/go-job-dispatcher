package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// MapWorker adds a mapping which stores what worker works on which job
// Todo: enhance this by adding a timestamp
func (m *RedisService) MapWorker(workerUID, jobUID uuid.UUID) (err error) {
	res := connections.RedisClient.HMSet(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String(), jobUID.String())
	err = res.Err()
	if err != nil {
		log.WithError(err).Error("Could not map worker to job")
	}
	return
}

// GetWorkerJob returns the uid of the job
func (m *RedisService) GetWorkerJob(workerUID uuid.UUID) (jobUID uuid.UUID, err error) {
	res := connections.RedisClient.HGet(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String())
	err = res.Err()
	if err != nil {
		log.WithError(err).Error("Could not get job uid from worker uid")
		return
	}
	// parse job uid
	jobUID, err = uuid.Parse(res.String())
	if err != nil {
		log.WithError(err).Error("Could not parse job uid")
	}
	return
}

// DeleteWorker removes a worker from the job
func (m *RedisService) DeleteWorker(workerUID, jobUID uuid.UUID) (err error) {
	res := connections.RedisClient.HDel(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String())
	err = res.Err()
	if err != nil {
		log.WithError(err).Error("Could not delete worker")
	}
	return
}
