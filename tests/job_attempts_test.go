package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAttempts(t *testing.T) {
	ass := assert.New(t)
	// create job
	_, _ = CreateJob()

	// get job
	err := Client.GetJob([]string{})
	ass.NoError(err)
	ass.Equal("test", Client.CurrentJob.WorkerInstance)
	ass.Equal(uint(1), Client.CurrentJob.Attempts)

	err = Client.ReleaseCurrentJob()
	ass.Nil(Client.CurrentJob)
	ass.NoError(err)

	err = Client.GetJob([]string{})
	ass.Equal("test", Client.CurrentJob.WorkerInstance)
	ass.Equal(uint(2), Client.CurrentJob.Attempts)

	err = Client.ReleaseCurrentJob()
	err = Client.GetJob([]string{})
	ass.Equal("test", Client.CurrentJob.WorkerInstance)
	ass.Equal(uint(3), Client.CurrentJob.Attempts)

	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
}
