package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetJob(t *testing.T) {
	// create a new job
	_, _ = CreateJob()
	ass := assert.New(t)
	err := Client.GetJob([]string{})
	ass.NoError(err)
	ass.Equal("test", Client.CurrentJob.WorkerInstance)
	err = Client.MarkCurrentJobAsCompleted()
	ass.NoError(err)
}
