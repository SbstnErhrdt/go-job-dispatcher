package job_dispatcher

import "github.com/google/uuid"

type JobService interface {
	New(*Job) error
	BulkNew([]*Job) error
	GetJobByUUID(uuid uuid.UUID) (*Job, error)
	Start(*Job) error
	HeartBeat(*Job, map[string]interface{}) error
	Release(*Job) error
	Complete(*Job) error
	GetLatestJob(workerInstances []string, workerUUID uuid.UUID) (*Job, error)
	GetCurrentJobOfWorker(workerInstances []string, workerUUID uuid.UUID) (*Job, bool, error)
	Clean() error
	GetStats() ([]Stats, error)
}
