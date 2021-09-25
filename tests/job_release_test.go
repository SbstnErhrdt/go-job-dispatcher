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
// Releases the job
func TestReleaseJob(t *testing.T) {
	ass := assert.New(t)
	// create job
	_, _ = CreateJob()
	// get job
	err := Client.GetJob([]string{})
	ass.NoError(err)
	ass.Equal("test", Client.CurrentJob.WorkerInstance)

	// save the id of the job
	jobID := Client.CurrentJob.UUID

	// release the job
	err = Client.ReleaseCurrentJob()
	ass.NoError(err)
	ass.Nil(Client.CurrentJob)

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
	ass.Nil(dbJob.CurrentWorkerUID)

	err = Client.GetJob([]string{})
	ass.NoError(err)
	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
}
