package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	connections.ConnectToRedis()
}

const (
	testStalledWorkInstance = "test_staled_instance_001"
)

var testStalledJob0 = job_dispatcher.Job{
	UUID:           uuid.MustParse("166c5bac-146e-11ec-82a8-0242ac130003"),
	WorkerInstance: testStalledWorkInstance,
}

var testStalledJob1 = job_dispatcher.Job{
	UUID:           uuid.MustParse("24167cf6-146e-11ec-82a8-0242ac130003"),
	WorkerInstance: testStalledWorkInstance,
}

// init service
var stalledTestService = RedisService{}

func TestGetStalledJobs(t *testing.T) {
	ass := assert.New(t)
	// add jobs
	err := stalledTestService.New(&testStalledJob0)
	ass.NoError(err)
	err = stalledTestService.New(&testStalledJob1)
	ass.NoError(err)
	// start job
	err = stalledTestService.Start(&testStalledJob0)
	ass.NoError(err)
	err = stalledTestService.Start(&testStalledJob1)
	ass.NoError(err)
	// stall one job
	// set the unix timestamp
	boolCMD := connections.RedisClient.HMSet(context.TODO(), GenerateKeyDoingMap(testStalledWorkInstance), testStalledJob0.UUID.String(), time.Now().UTC().Add(-time.Hour).Unix())
	ass.NoError(boolCMD.Err())
	// get the stalled jobs
	stalledJobs, err := GetStalledJobs()
	ass.Len(stalledJobs, 1)
	ass.Equal(testStalledJob0.UUID.String(), stalledJobs[0].String())
}
