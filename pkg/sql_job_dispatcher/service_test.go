package sql_job_dispatcher

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

var sqlService SqlService

func init() {
	log.Println("Load Test Env File")
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(path.Join(ex, "../"))
	filePath := exPath + "/.env"
	// Load environment
	err = godotenv.Load(filePath)
	if err != nil {
		log.Fatal("Error loading .env file", filePath)
	}

	connections.ConnectToSQL()
	connections.SQLClient.Logger.LogMode(logger.Info)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
}

func TestMain(m *testing.M) {
	log.Println("Start package test tests")
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	_ = m.Run()
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	log.Println("Done package test tests")
}

var t1 = job_dispatcher.Job{
	Name:           "testJob",
	WorkerInstance: "testWorker",
	Priority:       1,
	Parameters: map[string]string{
		"testCase1": "testCase1",
	},
	Tasks: []job_dispatcher.JobTask{
		{
			Name:    "testCase1",
			Version: "testCase1",
			Type:    "testCase1",
		},
	},
}
var t2 = job_dispatcher.Job{
	Name:           "highPriorityJob",
	WorkerInstance: "testWorker",
	Priority:       10,
	Parameters: map[string]string{
		"testCase1": "testCase1",
	},
	Tasks: []job_dispatcher.JobTask{
		{
			Name:    "testCase1",
			Version: "testCase1",
			Type:    "testCase1",
		},
	},
}

var t3 = job_dispatcher.Job{
	Name:           "lowPriorityJob",
	WorkerInstance: "testWorker",
	Priority:       7,
}
var t4 = job_dispatcher.Job{
	Name:           "highPriorityJob",
	WorkerInstance: "testWorker",
	Priority:       10,
}
var t5 = job_dispatcher.Job{
	Name:           "mediumPriorityJob",
	WorkerInstance: "testWorker",
	Priority:       9,
}

func TestNewJob(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	err := sqlService.New(&t1)
	ass.NoError(err)

	res := job_dispatcher.Job{}
	connections.SQLClient.First(&res)
	ass.Equal(uint(1), res.Priority)
	ass.Equal("testJob", res.Name)
	ass.Equal("testWorker", res.WorkerInstance)
	return
}

func TestGetLatestJob(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})

	// Test the GetCurrentJobOfWorker() method for the case when currently no job is assigned to the workerInstance
	// If there is no record for the workerInstance and uuid an empty job object is returned
	job, found, err := sqlService.GetCurrentJobOfWorker([]string{"testWorker"}, uuid.New())
	ass.False(found)
	ass.NoError(err)
	if err == nil {
		ass.Equal("", job.Name)
		ass.Nil(job.CurrentWorkerUID)
		ass.Equal("", job.WorkerInstance)
	}
	//ass.True(time.Time.IsZero(job.CompletedAt)) // Fails! TODO: Find out how to test a timestamp for nil

	// Test the GetCurrentJobOfWorker() func when job is assigned to workerInstance
	err = sqlService.New(&t2)
	ass.NoError(err)
	uid := uuid.New()
	// Call the GetLatest() method to assign the next most relevant job to the workerInstance
	job, err = sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	log.Println(job)
	job, found, err = sqlService.GetCurrentJobOfWorker([]string{"testWorker"}, uid)
	log.Println(job)
	if err == nil {
		ass.Equal("highPriorityJob", job.Name)
		ass.Equal(uid, *job.CurrentWorkerUID)
		ass.Equal("testWorker", job.WorkerInstance)
		ass.Equal(uint(10), job.Priority)
	}
	ass.True(found)
	ass.NoError(err)
	return
}

func TestPop(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	err := sqlService.New(&t3)
	err = sqlService.New(&t4)
	ass.NoError(err)
	err = sqlService.New(&t5)
	ass.NoError(err)
	uid := uuid.New()
	workerJob, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	ass.Equal("highPriorityJob", workerJob.Name)
	ass.Equal(uint(10), workerJob.Priority)
	err = sqlService.Complete(workerJob)
	ass.NoError(err)

	workerJob2, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	ass.Equal("mediumPriorityJob", workerJob2.Name)
	ass.Equal(uint(9), workerJob2.Priority)
	err = sqlService.Complete(workerJob2)
	ass.NoError(err)

	workerJob3, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	ass.Equal("lowPriorityJob", workerJob3.Name)
	ass.Equal(uint(7), workerJob3.Priority)
	err = sqlService.Complete(workerJob3)
	ass.NoError(err)

	log.Println(workerJob)
	return
}

func TestCleanStalledJobsHeartBeatIsNull(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	// create a new job
	err := sqlService.New(&t1)
	ass.NoError(err)
	// uuid of the test worker
	uid := uuid.New()
	// Assign the latest job
	workerJob, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	if err != nil {
		panic(err)
	}
	ass.Equal(uid, *workerJob.CurrentWorkerUID)
	err = sqlService.Clean()
	ass.NoError(err)
	// check the worker job
	res := job_dispatcher.Job{}
	connections.SQLClient.First(&res)
	ass.Equal(uid, *res.CurrentWorkerUID)
	ass.Nil(res.CompletedAt)
	return
}

func TestJobRelease(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	err := sqlService.New(&t4)
	ass.NoError(err)
	uid := uuid.New()
	// Call the GetLatest() method to assign the next most relevant job to the workerInstance
	workerJob, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	ass.Equal(uid, *workerJob.CurrentWorkerUID)
	err = sqlService.Release(workerJob)
	ass.NoError(err)
	res := job_dispatcher.Job{}
	err = connections.SQLClient.First(&res).Error
	ass.NoError(err)
	ass.Nil(res.CurrentWorkerUID)
	return
}

func TestCleanStalledJobsSetHeartBeat(t *testing.T) {
	ass := assert.New(t)
	_ = connections.SQLClient.Migrator().DropTable(&job_dispatcher.Job{})
	_ = connections.SQLClient.Migrator().CreateTable(&job_dispatcher.Job{})
	err := sqlService.New(&t4)
	ass.NoError(err)
	uid := uuid.New()
	workerJob, err := sqlService.GetLatestJob([]string{"testWorker"}, uid)
	ass.NoError(err)
	ass.Equal(uid, *workerJob.CurrentWorkerUID)
	err = sqlService.HeartBeat(workerJob, map[string]interface{}{})
	ass.NoError(err)
	// Simulate a stalled worker job that doesn't respond by waiting 11 minutes
	// Adjust the time and the 10 minutes interval in CleanStalledJobs() to not wait 11 minutes
	time.Sleep(22 * time.Second)
	err = sqlService.Clean()
	ass.NoError(err)
	res := job_dispatcher.Job{}
	connections.SQLClient.First(&res)
	ass.Nil(res.CurrentWorkerUID)
	ass.NotNil(res.LastHeartBeat)
	ass.Nil(res.CompletedAt)
	return
}

var a1 = job_dispatcher.Job{
	Name:           "a1",
	Priority:       10,
	WorkerInstance: "a1",
}
var b1 = job_dispatcher.Job{
	Name:           "b1",
	Priority:       10,
	WorkerInstance: "b1",
}
var c1 = job_dispatcher.Job{
	Name:           "c1",
	Priority:       10,
	WorkerInstance: "c1",
}

func TestGetStats(t *testing.T) {
	ass := assert.New(t)
	err := sqlService.New(&a1)
	err = sqlService.New(&b1)
	err = sqlService.New(&c1)
	err = sqlService.New(&a1)
	ass.NoError(err)
	res, err := sqlService.GetStats()
	ass.NoError(err)
	ass.NotEmpty(res)
}
