package tests

import (
	"context"
	"fmt"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/redis_job_dispatcher"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// Creates a job
// Gets the job
// Starts the job
// Completes job
func TestStartJob(t *testing.T) {
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
	ass.NotNil(Client.CurrentJob.StartedAt)

	time.Sleep(time.Second * 5)

	// start the job
	err = Client.HeartBeat(map[string]interface{}{"all": "fine"})
	ass.NoError(err)
	ass.NotNil(Client.CurrentJob)
	ass.NotNil(Client.CurrentJob.LastHeartBeat)

	time.Sleep(time.Second * 5)

	// check in the database if the job is also not assigned
	var dbJob job_dispatcher.Job
	if Client.Bulk {
		// retrieve job from redis
		log.Println("check redis")
		resCMD := connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyJobsMap(), jobID.String())
		ass.NoError(resCMD.Err())
		dbJob, err = redis_job_dispatcher.ParseJob(resCMD)
		ass.NoError(err)
	} else {
		// retrieve job from sql
		err = connections.SQLClient.First(&dbJob, jobID).Error
		ass.NoError(err)
	}
	ass.NoError(err)
	ass.NotNil(dbJob.StartedAt)
	fmt.Println(dbJob.StartedAt, dbJob.LastHeartBeat)
	ass.NotEqual(dbJob.StartedAt, dbJob.LastHeartBeat)

	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
	if Client.Bulk {
		// retrieve job from redis
		log.Println("check redis")
		resCMD := connections.RedisClient.HGet(context.TODO(), redis_job_dispatcher.GenerateKeyJobsMap(), jobID.String())
		ass.NoError(resCMD.Err())
		dbJob, err = redis_job_dispatcher.ParseJob(resCMD)
		ass.NoError(err)
	} else {
		// retrieve job from sql
		err = connections.SQLClient.First(&dbJob, jobID).Error
		ass.NoError(err)
	}
	ass.NotNil(dbJob.CompletedAt)
}
