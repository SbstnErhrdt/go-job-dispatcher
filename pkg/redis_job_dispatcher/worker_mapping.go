package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/google/uuid"
)

// MapWorker adds a mapping which stores what worker works on which job
// Todo: enhance this by adding a timestamp
func (m *RedisService) MapWorker(workerUID, jobUID uuid.UUID) error {
	res := connections.RedisClient.HMSet(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String(), jobUID.String())
	return res.Err()
}

// GetWorkerJob returns the uid of the job
func (m *RedisService) GetWorkerJob(workerUID uuid.UUID) (jobUID uuid.UUID, err error) {
	res := connections.RedisClient.HGet(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String())
	err = res.Err()
	if err != nil {
		return
	}
	// parse job uid
	jobUID, err = uuid.Parse(res.String())
	return
}

// DeleteWorker removes a worker from the job
func (m *RedisService) DeleteWorker(workerUID, jobUID uuid.UUID) error {
	res := connections.RedisClient.HDel(context.TODO(), GenerateKeyWorkerJobs(), workerUID.String())
	return res.Err()
}
