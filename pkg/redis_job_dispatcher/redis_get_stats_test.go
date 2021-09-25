package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	connections.ConnectToRedis()
}

const (
	testWorkInstance = "test_instance_001"
)

var testJob0 = job_dispatcher.Job{
	UUID:           uuid.New(),
	WorkerInstance: testWorkInstance,
}

var testJob1 = job_dispatcher.Job{
	UUID:           uuid.New(),
	WorkerInstance: testWorkInstance,
}

// init service
var service = RedisService{}

func setup() {
	connections.RedisClient.FlushDB(context.TODO())
}

func addJob0(ass *assert.Assertions) {
	// add job
	err := service.New(&testJob0)
	ass.NoError(err)
}

func addJob1(ass *assert.Assertions) {
	// add job
	err := service.New(&testJob1)
	ass.NoError(err)
}

func getJob(ass *assert.Assertions) (job *job_dispatcher.Job) {
	// add job
	// get the latest job
	job, err := service.GetLatestJob([]string{testWorkInstance}, uuid.New())
	ass.NoError(err)
	return
}

func startJob(ass *assert.Assertions, job *job_dispatcher.Job) {
	// complete job
	err := service.Start(job)
	ass.NoError(err)
}

func completeJob(ass *assert.Assertions, job *job_dispatcher.Job) {
	// complete job
	err := service.Complete(job)
	ass.NoError(err)
}

func TestGetKeys(t *testing.T) {
	ass := assert.New(t)
	// clear db
	setup()
	// check for keys
	keys, err := GetKeys()
	ass.NoError(err)
	ass.Len(keys, 0)

	// step 1: add job
	addJob0(ass)

	// check keys
	keys, err = GetKeys()
	ass.NoError(err)
	ass.Len(keys, 3)
	log.Println(keys)

	// step 2: add another job
	addJob1(ass)

	keys, err = GetKeys()
	ass.NoError(err)
	ass.Len(keys, 3)

	// step 3: get the job
	job := getJob(ass)
	// start the job
	startJob(ass, job)

	keys, err = GetKeys()
	ass.Len(keys, 5)

	// step 4: complete job
	completeJob(ass, job)

	keys, err = GetKeys()
	ass.Len(keys, 5)
}

func TestGetStats(t *testing.T) {
	ass := assert.New(t)
	// clear db
	setup()
	// check for keys
	stats, err := GetStats()
	ass.NoError(err)
	ass.Len(stats, 0)

	// step 1: add job
	addJob0(ass)

	// check keys
	stats, err = GetStats()
	ass.NoError(err)
	ass.Len(stats, 1)
	ass.Equal(1, stats[0].Todo)
	ass.Equal(0, stats[0].Active)
	ass.Equal(0, stats[0].Done)
	log.Println(stats)

	// step 2: add another job
	addJob1(ass)

	stats, err = GetStats()
	ass.NoError(err)
	ass.Len(stats, 1)
	ass.Equal(2, stats[0].Todo)
	ass.Equal(0, stats[0].Active)
	ass.Equal(0, stats[0].Done)
	log.Println(stats)

	// step 3: get the job
	job := getJob(ass)
	// start the job
	startJob(ass, job)

	stats, err = GetStats()
	ass.Len(stats, 1)
	ass.Equal(1, stats[0].Todo)
	ass.Equal(1, stats[0].Active)
	ass.Equal(0, stats[0].Done)
	log.Println(stats)

	// step 4: complete job
	completeJob(ass, job)

	stats, err = GetStats()
	ass.Len(stats, 1)
	ass.Equal(1, stats[0].Todo)
	ass.Equal(0, stats[0].Active)
	ass.Equal(1, stats[0].Done)
	log.Println(stats)
}
