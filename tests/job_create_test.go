package tests

import (
	"bytes"
	"encoding/json"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateJobFail(t *testing.T) {
	ass := assert.New(t)
	// create payload
	payload := map[string]interface{}{
		"hello": "world",
	}
	// marshal payload
	jsonPayload, _ := json.Marshal(payload)
	// send request
	resp, err := http.Post(Server.URL+"/jobs/", "application/json", bytes.NewBuffer(jsonPayload))
	ass.NoError(err)
	ass.Equal(400, resp.StatusCode)
}

func CreateJob() (job *job_dispatcher.Job, err error) {
	newJob := job_dispatcher.NewJobDTO{
		Name:           "test",
		WorkerInstance: "test",
		Priority:       9,
		Params: map[string]string{
			"test": "test",
		},
		Tasks: []job_dispatcher.JobTask{
			{
				Version: "1.0",
				Name:    "Test",
				Execute: nil,
				Type:    "Nothing",
			},
		},
	}

	job, err = Client.CreateJob(newJob)
	return
}

func CreateJobWithOtherInstance() (job *job_dispatcher.Job, err error) {
	newJob := job_dispatcher.NewJobDTO{
		Name:           "test",
		WorkerInstance: "test2",
		Priority:       1,
		Params: map[string]string{
			"test": "test",
		},
		Tasks: []job_dispatcher.JobTask{
			{
				Version: "1.0",
				Name:    "Test",
				Execute: nil,
				Type:    "Nothing",
			},
		},
	}

	job, err = Client.CreateJob(newJob)
	return
}

func TestCreateJob(t *testing.T) {
	ass := assert.New(t)
	job, err := CreateJob()
	ass.NoError(err)
	ass.NotEqual(0, job.UUID)
	err = Client.GetJob([]string{})
	ass.NoError(err)
	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
}

func TestCreateJobWithOtherInstance(t *testing.T) {
	ass := assert.New(t)
	job, err := CreateJobWithOtherInstance()
	ass.NoError(err)
	ass.NotEqual(0, job.UUID)
	err = Client.GetJob([]string{"test2"})
	ass.NoError(err)
	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
}
