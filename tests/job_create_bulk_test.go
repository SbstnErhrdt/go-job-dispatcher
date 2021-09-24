package tests

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateBulkJobs(t *testing.T) {
	ass := assert.New(t)

	newJobs := []job_dispatcher.NewJobDTO{
		{
			Name:           "bulk_test",
			WorkerInstance: "bulk_test_1",
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
		},
		{
			Name:           "bulk_test",
			WorkerInstance: "bulk_test_2",
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
		},
	}

	res, err := Client.CreateJobs(newJobs)

	ass.NoError(err)
	for _, j := range res {
		ass.NotEqual(0, j.UUID)
	}

}
