package redis_job_dispatcher

import (
	"context"
	"encoding/json"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type RedisService struct{}

// New Creates a new job
func (m *RedisService) New(job *job_dispatcher.Job) (err error) {
	ctx := context.TODO()
	// job to json
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.WithError(err).Error("Could not marshal job")
		return err
	}
	// init transaction
	tx := connections.RedisClient.TxPipeline()
	// create the job in the global list
	tx.HMSet(context.TODO(), GenerateKeyJobsMap(), job.UUID.String(), string(jobBytes))
	// add all keyJobs to the redis key value set
	tx.HMSet(ctx, GenerateKeyTodoMap(job.WorkerInstance), job.UUID.String(), 1)
	// add the uuid to the keyQueue
	tx.LPush(ctx, GenerateKeyTodoQueue(job.WorkerInstance), job.UUID.String())
	// execute the transaction
	_, err = tx.Exec(ctx)
	if err != nil {
		log.WithError(err).Error("Could not create job")
	}
	return err
}

// BulkNew creates multiple new keyJobs
func (m *RedisService) BulkNew(newJobs []*job_dispatcher.Job) (err error) {
	tx := connections.RedisClient.TxPipeline()
	for _, j := range newJobs {
		// job to json
		jobBytes, errMarshal := json.Marshal(j)
		if errMarshal != nil {
			err = errMarshal
			log.WithError(err).Error("Could not marshal job")
			return
		}
		// create the job in the global list
		tx.HMSet(context.TODO(), GenerateKeyJobsMap(), j.UUID.String(), string(jobBytes))
		// add all keyJobs to the redis key value set
		tx.HMSet(context.TODO(), GenerateKeyTodoMap(j.WorkerInstance), j.UUID.String(), 1)
		// add the uuid to the keyQueue
		tx.LPush(context.TODO(), GenerateKeyTodoQueue(j.WorkerInstance), j.UUID.String())
	}
	// execute the transaction
	_, err = tx.Exec(context.TODO())
	if err != nil {
		log.WithError(err).Error("Could not create job")
		return
	}
	return
}

// ParseJob extracts the job from a redis result object
func ParseJob(cmd *redis.StringCmd) (job job_dispatcher.Job, err error) {
	var resString string
	err = cmd.Scan(&resString)
	if err != nil {
		log.WithError(err).Error("Could not scan job from redis")
		return
	}
	err = json.Unmarshal([]byte(resString), &job)
	if err != nil {
		log.WithError(err).Error("Could not unmarshal job from redis")
		return
	}
	return
}

// GetJobByUUID returns a job by its unique identifier
func (m *RedisService) GetJobByUUID(uid uuid.UUID) (result *job_dispatcher.Job, err error) {
	resCMD := connections.RedisClient.HGet(context.TODO(), GenerateKeyJobsMap(), uid.String())
	job, err := ParseJob(resCMD)
	if err != nil {
		log.WithError(err).Error("Could not parse job")
		return
	}
	result = &job
	return
}

func (m *RedisService) setJobActive(job *job_dispatcher.Job) (err error) {
	ctx := context.TODO()
	tx := connections.RedisClient.TxPipeline()
	// remove entry from tdo map
	tx.HDel(ctx, GenerateKeyTodoMap(job.WorkerInstance), job.UUID.String())
	// add entry to keyDoing map
	tx.HMSet(ctx, GenerateKeyDoingMap(job.WorkerInstance), job.UUID.String(), time.Now().UTC().Unix())
	// remove potential entry from keyDone map
	tx.HDel(ctx, GenerateKeyDoneMap(job.WorkerInstance), job.UUID.String())
	// save job to json
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.WithError(err).Error("Could not marshal job")
		return err
	}
	tx.HMSet(
		context.TODO(),
		GenerateKeyJobsMap(),
		job.UUID.String(),
		string(jobBytes),
	)
	// execute the transaction
	_, err = tx.Exec(context.TODO())
	if err != nil {
		log.WithError(err).Error("Could not set job active")
	}
	return
}

func (m *RedisService) deleteJob(job *job_dispatcher.Job) (err error) {
	ctx := context.TODO()
	tx := connections.RedisClient.TxPipeline()
	// remove entry from tdo map
	tx.HDel(ctx, GenerateKeyTodoMap(job.WorkerInstance), job.UUID.String())
	// remove entry from keyDoing map
	tx.HDel(ctx, GenerateKeyDoingMap(job.WorkerInstance), job.UUID.String())
	// remove entry from keyDone map
	tx.HDel(ctx, GenerateKeyDoneMap(job.WorkerInstance), job.UUID.String())
	// remove job from normal map
	tx.HDel(ctx, GenerateKeyJobsMap(), job.UUID.String())
	// execute the transaction
	_, err = tx.Exec(context.TODO())
	return
}

// Start marks a job as started
func (m *RedisService) Start(job *job_dispatcher.Job) error {
	// save the job
	now := time.Now()
	job.StartedAt = &now
	return m.setJobActive(job)
}

// HeartBeat marks a job as alive
func (m *RedisService) HeartBeat(job *job_dispatcher.Job, status map[string]interface{}) error {
	now := time.Now()
	job.LastHeartBeat = &now
	job.CurrentStatus = status
	return m.setJobActive(job)
}

// Release releases the job
func (m *RedisService) Release(job *job_dispatcher.Job) (err error) {
	job.LastHeartBeat = nil
	job.StartedAt = nil
	ctx := context.TODO()
	tx := connections.RedisClient.TxPipeline()
	// add entry to tdo map
	tx.HMSet(ctx, GenerateKeyTodoMap(job.WorkerInstance), job.UUID.String(), 1)
	// add the uuid to the keyQueue
	tx.LPush(ctx, GenerateKeyTodoQueue(job.WorkerInstance), job.UUID.String())
	// remove entry from keyDoing map
	tx.HDel(ctx, GenerateKeyDoingMap(job.WorkerInstance), job.UUID.String())
	// remove entry from keyDone map
	tx.HDel(ctx, GenerateKeyDoneMap(job.WorkerInstance), job.UUID.String())
	// save job
	// transform job to json
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.WithError(err)
		return err
	}
	tx.HMSet(
		context.TODO(),
		GenerateKeyJobsMap(),
		job.UUID.String(),
		string(jobBytes),
	)
	// execute the transaction
	_, err = tx.Exec(context.TODO())
	if err != nil {
		log.WithError(err).Error("Could not release job")
	}
	return
}

// Complete marks the job as completed
func (m *RedisService) Complete(job *job_dispatcher.Job) (err error) {
	now := time.Now()
	job.LastHeartBeat = &now
	job.CompletedAt = &now

	ctx := context.TODO()
	tx := connections.RedisClient.TxPipeline()
	// remove entry from tdo map
	tx.HDel(ctx, GenerateKeyTodoMap(job.WorkerInstance), job.UUID.String())
	// remove entry from keyDoing map
	tx.HDel(ctx, GenerateKeyDoingMap(job.WorkerInstance), job.UUID.String())
	// add entry to keyDone map
	tx.HMSet(ctx, GenerateKeyDoneMap(job.WorkerInstance), job.UUID.String(), 1)
	// save job
	// job to json
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.WithError(err).Error("Could not marshal job")
		return err
	}
	tx.HMSet(
		context.TODO(),
		GenerateKeyJobsMap(),
		job.UUID.String(),
		string(jobBytes),
	)
	// execute the transaction
	_, err = tx.Exec(context.TODO())
	if err != nil {
		log.WithError(err).Error("Could not complete job")
	}
	return
}

// GetLatestJob returns the latest job
func (m *RedisService) GetLatestJob(workerInstances []string, workerUUID uuid.UUID) (job *job_dispatcher.Job, err error) {
	var jobUIDString string
	// shuffle instances
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(workerInstances), func(i, j int) { workerInstances[i], workerInstances[j] = workerInstances[j], workerInstances[i] })
	// iterate over the instances
	for _, instance := range workerInstances {
		// pop the job from the list
		res := connections.RedisClient.LPop(context.TODO(), GenerateKeyTodoQueue(instance))
		resErr := res.Err()
		if resErr != nil {
			continue
		}
		errScan := res.Scan(&jobUIDString)
		if errScan != nil {
			err = errScan
			log.WithError(err).Error("Could not scan job uid string")
			return
		} else {
			break
		}
	}
	if len(jobUIDString) == 0 {
		err = job_dispatcher.ErrNoNewJobs
		return nil, err
	}
	// return the job
	jobUID, err := uuid.Parse(jobUIDString)
	if err != nil {
		log.WithError(err).Error("Could not parse job uid")
		return
	}
	// map the worker with the uuid
	err = m.MapWorker(workerUUID, jobUID)
	if err != nil {
		log.WithError(err).Error("Could not map worker with job")
		return
	}
	// get the job
	job, err = m.GetJobByUUID(jobUID)
	if err != nil {
		log.WithError(err).Error("Could not get job by uuid")
		return
	}
	// increment the attempts
	job.Attempts += 1
	// save the job
	err = m.SaveJob(job)
	if err != nil {
		log.WithError(err).Error("Could not save job")
		return
	}
	return
}

// GetCurrentJobOfWorker retrieves the job of the worker if one exists in the database
func (m *RedisService) GetCurrentJobOfWorker(
	workerInstances []string,
	currentWorkerUID uuid.UUID,
) (job *job_dispatcher.Job, found bool, err error) {
	jobUId, err := m.GetWorkerJob(currentWorkerUID)
	if err == nil {
		found = false
		err = nil
		return
	}
	job, err = m.GetJobByUUID(jobUId)
	if err != nil {
		log.WithError(err).Error("Could not get job by uuid")
		return
	}
	found = true
	return
}

// Clean cleans all stalled keyJobs
// get different job instances
// find all jobs of the instances in the active map
// iterate over the jobs and compare the timestamps
// if the timestamp is larger than the threshold
// add the job to the queue again
func (m *RedisService) Clean() error {
	stalledJobUIDs, err := GetStalledJobs()
	if err != nil {
		log.WithError(err).Error("Could not get stalled jobs")
		return err
	}
	for _, jobUID := range stalledJobUIDs {
		// get the stalled job object
		stalledJob, warnJob := m.GetJobByUUID(jobUID)
		if warnJob != nil {
			log.WithField("jobUID", jobUID).Warn(warnJob)
			continue
		}
		// add it to the queue again
		warnRelease := m.Release(stalledJob)
		if warnRelease != nil {
			log.WithField("jobUID", jobUID).Warn(warnRelease)
			continue
		}
	}
	return nil
}

// GetStats returns the stats of all the jobs
func (m *RedisService) GetStats() (results []job_dispatcher.Stats, err error) {
	return GetStats()
}

// SaveJob saves a job in the database
func (m *RedisService) SaveJob(job *job_dispatcher.Job) (err error) {
	// job to json
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.WithError(err).Error("Could not marshal job")
		return err
	}
	// create the job in the global list
	resCMD := connections.RedisClient.HMSet(context.TODO(), GenerateKeyJobsMap(), job.UUID.String(), string(jobBytes))
	if err != nil {
		log.WithError(err).Error("Could not save job")
		return
	}
	err = resCMD.Err()
	if err != nil {
		log.WithError(err).Error("Could not save job")
		return
	}
	return
}
