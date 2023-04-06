package job_dispatcher

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// ResponseDTO is the response for a single job
type ResponseDTO struct {
	Status string `json:"status"`
	Job    *Job   `json:"job"`
}

// ResponseMultipleDTO is the response for multiple jobs
type ResponseMultipleDTO struct {
	Status string `json:"status"`
	Jobs   []*Job `json:"jobs"`
}

// NewJobDTO is the dto for creating a new job
type NewJobDTO struct {
	Name           string            `json:"name"`
	WorkerInstance string            `json:"workerInstance"`
	Priority       uint              `json:"priority"`
	Params         map[string]string `json:"params"`
	Tasks          []JobTask         `json:"tasks"`
}

// ErrPayloadHasNoName is returned if the payload has no name
var ErrPayloadHasNoName = errors.New("payload has no name")

// ErrWorkerInstanceNotSet is returned if the worker instance is not set
var ErrWorkerInstanceNotSet = errors.New("worker instance not set")

// Validate checks if a job is valid
func (dto *NewJobDTO) Validate() (err error) {
	if len(dto.Name) == 0 {
		return ErrPayloadHasNoName
	}
	if len(dto.WorkerInstance) == 0 {
		return ErrWorkerInstanceNotSet
	}
	return
}

// GenerateJob transforms a dto into a job
// metadata is added
// new uuid is generated
func (dto *NewJobDTO) GenerateJob() *Job {
	return &Job{
		UUID:           uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      nil,
		Name:           dto.Name,
		Priority:       dto.Priority,
		WorkerInstance: dto.WorkerInstance,
		Parameters:     dto.Params,
		Tasks:          dto.Tasks,
	}
}
