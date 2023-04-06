package sql_job_dispatcher

import (
	"encoding/json"
	"errors"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
)

type SqlService struct{}

func (m *SqlService) New(job *job_dispatcher.Job) (err error) {
	sqlJob := Job{}
	sqlJob.Job = *job
	err = connections.SQLClient.Create(&sqlJob).Error
	return err
}

// BulkNew creates multiple new jobs
func (m *SqlService) BulkNew(newJobs []*job_dispatcher.Job) (err error) {
	// init result slice
	sqlResults := make([]*Job, len(newJobs))
	// iterate and copy
	for i, j := range newJobs {
		tempJob := j
		sqlJob := Job{}
		sqlJob.Job = *tempJob
		sqlResults[i] = &sqlJob
	}
	// create in the db
	err = connections.SQLClient.Create(&sqlResults).Error
	// return newJobs or err
	return
}

// GetJobByUUID returns a job by its unique identifier
func (m *SqlService) GetJobByUUID(uid uuid.UUID) (job *job_dispatcher.Job, err error) {
	sqlJob := Job{}
	err = connections.SQLClient.
		Where("uuid=?", uid).
		First(&sqlJob).Error
	if err != nil {
		return
	}
	job = &sqlJob.Job
	return
}

// Start marks a job as started
func (m *SqlService) Start(job *job_dispatcher.Job) (err error) {
	err = connections.SQLClient.Exec(`
		UPDATE jobs 
		SET started_at=current_timestamp 
		WHERE uuid = ?
	`, job.UUID).Error
	return nil
}

// HeartBeat marks a job as alive
func (m *SqlService) HeartBeat(job *job_dispatcher.Job, status map[string]interface{}) (err error) {
	data, _ := json.Marshal(status)
	err = connections.SQLClient.Exec(`
		UPDATE jobs 
		SET 
			last_heart_beat = current_timestamp,
			current_status_json = ? 
		WHERE uuid = ?
	`, data, job.UUID).Error
	return nil
}

// Release releases the job
func (m *SqlService) Release(job *job_dispatcher.Job) (err error) {
	err = connections.SQLClient.Exec(`
		UPDATE jobs 
		SET current_worker_uid = NULL 
		WHERE uuid = ?
	`, job.UUID).Error
	return nil
}

// Complete marks the job as completed
func (m *SqlService) Complete(job *job_dispatcher.Job) (err error) {
	err = connections.SQLClient.Exec(`
		UPDATE jobs 
		SET completed_at=current_timestamp 
		WHERE uuid = ?
	`, job.UUID).Error
	return
}

// GetLatestJob returns the latest job
func (m *SqlService) GetLatestJob(workerInstances []string,
	currentWorkerUID uuid.UUID) (job *job_dispatcher.Job, err error) {
	// Step 1: check if there is already a job on which the worker is currently working
	job, found, err := m.GetCurrentJobOfWorker(workerInstances, currentWorkerUID)
	// If there is already a job
	// -> terminate and return the job
	if found {
		return
	}
	// if there is an error
	if err != nil {
		return
	}

	// If there is not a job at the moment
	// Step 2: mark the latest element in the queue with the uuid of the current worker
	// By doing that this job is now assigned in an atomic way to the worker bot
	err = connections.SQLClient.Exec(`
		UPDATE jobs as w1,
		(
			SELECT uuid 
			FROM jobs 
			WHERE worker_instance IN ? 
			AND completed_at IS NULL
			AND current_worker_uid IS NULL
			ORDER BY priority DESC, attempts ASC
			LIMIT 1
		) as w2		
		SET w1.current_worker_UID=?,
		last_heart_beat=current_timestamp,
		attempts = attempts + 1
		WHERE w1.uuid = w2.uuid		 
	`, workerInstances, currentWorkerUID).Error
	if err != nil {
		return
	}

	// Step 3: Get this new job
	job, found, err = m.GetCurrentJobOfWorker(workerInstances, currentWorkerUID)
	// if there is an error
	if err != nil {
		return
	}
	// if there is still no job
	if !found {
		err = job_dispatcher.ErrNoNewJobs
		return
	}
	return
}

// GetCurrentJobOfWorker retrieves the job of the worker if one exists in the database
func (m *SqlService) GetCurrentJobOfWorker(
	workerInstances []string,
	currentWorkerUID uuid.UUID,
) (job *job_dispatcher.Job, found bool, err error) {
	res := Job{}
	err = connections.SQLClient.
		Where("worker_instance in ?", workerInstances).
		Where("current_worker_UID=?", currentWorkerUID).
		Where("completed_at IS NULL").
		Order("priority DESC").
		Order("attempts ASC").
		First(&res).Error
	job = &res.Job
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If the record not was found
		// terminate
		found = false
		err = nil
		return
	} else {
		// if there is no record but an error
		// return the error
		if err != nil {
			found = false
			return
		}
		// if there is a record
		// return the record
		found = true
		return
	}
}

// Clean cleans all stalled jobs
func (m *SqlService) Clean() (err error) {
	stalledJobsLimit := os.Getenv("CLEAN_STALLED_JOBS_INTERVAL")
	if len(stalledJobsLimit) == 0 {
		stalledJobsLimit = "10 MINUTE"
	}
	// Todo find better way to include interval
	err = connections.SQLClient.Exec(`
		UPDATE jobs 
		SET current_worker_UID=NULL 
		WHERE completed_at IS NULL 
		AND current_worker_UID IS NOT NULL
		AND (
			last_heart_beat < current_timestamp - interval ` + stalledJobsLimit + `
			OR last_heart_beat IS NULL
		)
	`).Error
	return
}

// GetStats returns the statistics of each job instance
func (m *SqlService) GetStats() (results []job_dispatcher.Stats, err error) {
	err = connections.SQLClient.Raw(`
SELECT
	worker_instance,
	SUM(CASE WHEN 1 THEN 1 ELSE 0 END) as total,
	SUM(CASE WHEN completed_at IS NULL THEN 1 ELSE 0 END) as todo,
	SUM(CASE WHEN completed_at IS NOT NULL THEN 1 ELSE 0 END) as done,
	SUM(CASE WHEN last_heart_beat > date_sub(now(), interval 60 second) THEN 1 ELSE 0 END) as active
FROM jobs
GROUP BY worker_instance
	`).Find(&results).Error
	return
}
