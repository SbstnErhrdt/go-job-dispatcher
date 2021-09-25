package tests

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/redis_job_dispatcher"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Creates a job
// Gets the job
// Sends heartbeat
// Completes job
func TestHeartbeatJob(t *testing.T) {
	ass := assert.New(t)
	// create job
	_, _ = CreateJob()
	// get job
	err := Client.GetJob([]string{})
	ass.NoError(err)
	ass.Equal("test", Client.CurrentJob.WorkerInstance)

	// save the id of the job
	jobID := Client.CurrentJob.UUID

	// start the job
	err = Client.StartCurrentJob()
	ass.NoError(err)
	ass.NotNil(Client.CurrentJob)
	ass.NotNil(jobID)
	ass.NotNil(Client.CurrentJob.StartedAt)

	var dbJob job_dispatcher.Job
	if Client.Bulk {
		log.Println("check redis")
		resCMD := connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyJobsMap(), jobID.String())
		ass.NoError(resCMD.Err())
		dbJob, err = redis_job_dispatcher.ParseJob(resCMD)
		ass.NoError(err)
	} else {
		log.Println("check sql")
		// check in the sql database if the job is also not assigned
		err = connections.SQLClient.First(&dbJob, jobID).Error
		ass.NoError(err)
	}
	ass.NotNil(dbJob.StartedAt)
	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
	if Client.Bulk {
		log.Println("check redis")
		// check if in the result set is the key with the 1
		resCMD := connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyDoneMap(dbJob.WorkerInstance), dbJob.UUID.String())
		ass.NoError(resCMD.Err())
		resInt, errInt := resCMD.Int()
		ass.NoError(errInt)
		ass.Equal(1, resInt)
		// not results when looking at the tdo and doing list
		resCMD = connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyTodoMap(dbJob.WorkerInstance), dbJob.UUID.String())
		ass.Error(resCMD.Err())
		resCMD = connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyDoingMap(dbJob.WorkerInstance), dbJob.UUID.String())
		ass.Error(resCMD.Err())
	} else {
		log.Println("check sql")
		err = connections.SQLClient.First(&dbJob, jobID).Error
		ass.NoError(err)
		ass.NotNil(dbJob.CompletedAt)
	}

}
