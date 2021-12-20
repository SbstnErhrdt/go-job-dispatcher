package job_dispatcher

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Job struct {
	UUID uuid.UUID `json:"uuid" gorm:"primaryKey"`
	// Metadata
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	// UUIDs
	MandateUID uuid.UUID `json:"mandateUID" gorm:"index"`
	ClientUID  uuid.UUID `json:"clientUID" gorm:"index"`
	OwnerUID   uuid.UUID `json:"ownerUID" gorm:"index"`
	// Attributes
	StartedAt         *time.Time             `json:"startedAt"`                                                     // the timestamp when the worker started on this job
	CompletedAt       *time.Time             `json:"completedAt" gorm:"index"`                                      // the timestamp when the worker finished the job
	LastHeartBeat     *time.Time             `json:"lastHeartBeat"`                                                 // the last time the worker was working on this job
	Name              string                 `json:"name" gorm:"type:varchar(100);"`                                // a short and optional name
	Priority          uint                   `json:"priority" gorm:"default:1;"`                                    // the priority of the job (1 lowest priority, ... 10 high priority)
	Attempts          uint                   `json:"attempts" gorm:"default:0;"`                                    // the attempts that have been tried in the past
	CurrentWorkerUID  *uuid.UUID             `json:"currentWorkerUID" gorm:"index; type:varchar(36); default:null"` // the uuid of the the worker bot that is currently working in this job
	WorkerInstance    string                 `json:"workerInstance" gorm:"index"`                                   // the type of worker that should work on this job
	Parameters        map[string]string      `json:"parameters" gorm:"-"`                                           // the search parameters the worker should use (e.g. company names ... )
	ParametersJSON    datatypes.JSON         `json:"-" gorm:"type:jsonb;"`                                          // json
	Tasks             []JobTask              `json:"tasks" gorm:"-"`                                                // the way the worker should interact (e.g. click on that, extract this, store it in this bucket, ...)
	TasksJSON         datatypes.JSON         `json:"-" gorm:"type:jsonb;"`                                          // json
	CurrentStatus     map[string]interface{} `json:"currentStatus" gorm:"-"`                                        // the current status of the job
	CurrentStatusJSON datatypes.JSON         `json:"-" gorm:"type:jsonb;"`                                          // json
}

// JobTask is a task the worker should perform
type JobTask struct {
	Version string      `json:"version"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Execute interface{} `json:"execute"`
}

func (wj *Job) String() string {
	return fmt.Sprintf("<Job id: %v name: %s>", wj.UUID, wj.Name)
}
