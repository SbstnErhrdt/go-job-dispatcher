package api

import (
	"errors"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/redis_job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/sql_job_dispatcher"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func getService(c *gin.Context) job_dispatcher.JobService {
	// get the service from the context
	pass := c.MustGet(JobServiceKey)
	var s job_dispatcher.JobService
	switch pass.(type) {
	case *sql_job_dispatcher.SqlService:
		s = pass.(job_dispatcher.JobService)
		break
	case *redis_job_dispatcher.RedisService:
		s = pass.(job_dispatcher.JobService)
		break
	default:
		panic(errors.New("can not find right type"))
	}
	return s
}

// extracts and parses id from params
func extractUUIDFromPath(c *gin.Context) (res uuid.UUID, err error) {
	uuidString := c.Param(UuidPathKey)
	res, err = uuid.Parse(uuidString)
	return
}

// get the job from the database based on the id
func getJob(c *gin.Context) (job *job_dispatcher.Job, err error) {
	// get the service from the context
	s := getService(c)
	// extract the uid from the params
	uid, err := extractUUIDFromPath(c)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "no id provided",
		})
		return
	}
	// get the job
	job, err = s.GetJobByUUID(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "no id provided",
		})
		return
	}
	return
}

// CreateJobHandler creates a new job
func CreateJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// init the payload
	var payload job_dispatcher.NewJobDTO
	// parse the payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	// validate
	err = payload.Validate()
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	newJob := payload.GenerateJob()
	// create the new job
	err = s.New(newJob)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	// return the new job
	c.JSON(200, job_dispatcher.ResponseDTO{
		Status: "created",
		Job:    newJob,
	})
	return
}

// BulkCreateJobsHandler create multiple new jobs
func BulkCreateJobsHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// init the payload
	var payload []job_dispatcher.NewJobDTO
	// parse the payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	// iterate over all the jobs and validate them
	for _, j := range payload {
		err = j.Validate()
		if err != nil {
			c.JSON(400, gin.H{
				"err": err,
				"job": j,
			})
			return
		}
	}
	// generate jobs
	var jobs []*job_dispatcher.Job
	for _, j := range payload {
		// validate each job
		err = j.Validate()
		if err != nil {
			return
		}
		// generate a new job from the dto
		res := j.GenerateJob()
		// append the job to the result set
		jobs = append(jobs, res)
	}
	// create the new jobs
	err = s.BulkNew(jobs)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	// return the new job
	c.JSON(200, job_dispatcher.ResponseMultipleDTO{
		Status: "created",
		Jobs:   jobs,
	})
	return
}

// StartJobHandler marks a job as started
func StartJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// get the job
	job, err := getJob(c)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = s.Start(job)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	job, err = getJob(c)
	if err != nil {
		c.JSON(500, gin.H{
			"err": "can not start job",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"status": "started",
			"job":    job,
		})
	}
	return
}

// HeartBeatJobHandler mark job as alive
func HeartBeatJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// init the payload
	var payload map[string]interface{}
	// parse the payload
	err := c.BindJSON(&payload)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	job, err := getJob(c)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = s.HeartBeat(job, payload)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	job, err = getJob(c)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{
			"err": "can not start job",
		})
		return
	} else {
		c.JSON(200, job_dispatcher.ResponseDTO{
			Status: "alive",
			Job:    job,
		})
	}
	return
}

// ReleaseJobHandler Releases the job after crash or problem
func ReleaseJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// get the job
	job, err := getJob(c)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = s.Release(job)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	job, err = getJob(c)
	if err != nil {
		c.JSON(500, gin.H{
			"err": "can not release job",
		})
		return
	} else {
		c.JSON(200, job_dispatcher.ResponseDTO{
			Status: "released",
			Job:    job,
		})
	}
	return
}

// CompleteJobHandler marks a job completed
func CompleteJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	job, err := getJob(c)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = s.Complete(job)
	if err != nil {
		log.Error(err)
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	job, err = getJob(c)
	if err != nil {
		c.JSON(500, gin.H{
			"err": "can not complete job",
		})
		return
	} else {
		c.JSON(200, job_dispatcher.ResponseDTO{
			Status: "completed",
			Job:    job,
		})
	}
	return
}

// CleanJobsHandler is the handler to clean the stalled jobs
func CleanJobsHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	err := s.Clean()
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"err": err,
		})
	} else {
		c.Status(200)
	}
}

// GetJobsHandler gets the latest job from the database
// Checks if there is already a job assigned
func GetJobsHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// parse params
	instance := c.Param("instance")
	uuidString := c.Param("uuid")
	// parse uuid string to uuid object
	uid, err := uuid.Parse(uuidString)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"err": err,
		})
		return
	}
	// init the instances array
	instances := []string{instance}
	// if the request has a payload
	type RequestBody struct {
		Instances []string `json:"instances"`
	}
	reqBody := RequestBody{}
	err = c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"err": err,
		})
		return
	}
	// add the additional instances to the array if there are some
	if len(reqBody.Instances) > 0 {
		instances = append(instances, reqBody.Instances...)
	}
	// pop the latest job
	job, err := s.GetLatestJob(instances, uid)
	if err != nil {
		// return 204 if there is no new jobs
		if err == job_dispatcher.ErrNoNewJobs {
			c.JSON(204, nil)
			return
		}
		// if there is a real error
		c.JSON(500, map[string]interface{}{
			"err": err,
		})
		return
	} else {
		// return the job
		c.JSON(200, job_dispatcher.ResponseDTO{
			Status: "assigned",
			Job:    job,
		})
		return
	}
}

// ReadStatsJobHandler reads the current stats
func ReadStatsJobHandler(c *gin.Context) {
	// get the service from the context
	s := getService(c)
	// read the stats
	results, err := s.GetStats()
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	// return the stats
	c.JSON(200, results)
	return
}
