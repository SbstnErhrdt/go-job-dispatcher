package redis_job_dispatcher

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testInstance = "test"
)

func init() {
	connections.ConnectToRedis()
}

func TestPop(t *testing.T) {
	ass := assert.New(t)
	job := job_dispatcher.Job{
		Name:           testInstance,
		Priority:       1,
		WorkerInstance: testInstance,
		Parameters:     map[string]string{},
		Tasks:          []job_dispatcher.JobTask{},
	}
	ass.NotEmpty(job)
	err := Push(job)
	ass.NoError(err)

	dbJob, err := Pop(testInstance)
	ass.NoError(err)
	ass.Equal(job.Priority, dbJob.Priority)
}

func TestPush(t *testing.T) {
	ass := assert.New(t)
	job := job_dispatcher.Job{
		Name:           testInstance,
		Priority:       10,
		WorkerInstance: testInstance,
		Parameters:     map[string]string{},
		Tasks:          []job_dispatcher.JobTask{},
	}
	ass.NotEmpty(job)
	err := Push(job)
	ass.NoError(err)

	dbJob, err := Pop(testInstance)
	ass.NoError(err)
	ass.Equal(job.Priority, dbJob.Priority)
}
