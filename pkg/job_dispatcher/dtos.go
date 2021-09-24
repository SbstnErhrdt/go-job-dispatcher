package job_dispatcher

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type ResponseDTO struct {
	Status string `json:"status"`
	Job    *Job   `json:"job"`
}

type ResponseMultipleDTO struct {
	Status string `json:"status"`
	Jobs   []*Job `json:"jobs"`
}

type NewJobDTO struct {
	Name           string            `json:"name"`
	WorkerInstance string            `json:"workerInstance"`
	Priority       uint              `json:"priority"`
	Params         map[string]string `json:"params"`
	Tasks          []JobTask         `json:"tasks"`
}

// Validate checks if a job is valid
func (dto *NewJobDTO) Validate() (err error) {
	if len(dto.Name) == 0 {
		return errors.New("payload has no name")
	}
	if len(dto.WorkerInstance) == 0 {
		return errors.New("worker instance is not set")
	}
	return
}

// GenerateJob transforms a dto into a job
// meta data is added
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
