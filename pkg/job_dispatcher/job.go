package job_dispatcher

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Job struct {
	UUID uuid.UUID `gorm:"type:varchar(36); primaryKey" json:"uuid"`
	// Metadata
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`
	// UUIDs
	MandateUID uuid.UUID `gorm:"index" gorm:"type:varchar(36);" json:"mandateUID"`
	ClientUID  uuid.UUID `gorm:"index" gorm:"type:varchar(36);" json:"clientUID"`
	OwnerUID   uuid.UUID `gorm:"index" gorm:"type:varchar(36);" json:"ownerUID"`
	// Attributes
	StartedAt         *time.Time             `gorm:"index" json:"startedAt"`                                        // the timestamp when the worker started on this job
	CompletedAt       *time.Time             `gorm:"index" json:"completedAt"`                                      // the timestamp when the worker finished the job
	LastHeartBeat     *time.Time             `gorm:"index" json:"lastHeartBeat"`                                    // the last time the worker was working on this job
	Name              string                 `json:"name" gorm:"type:varchar(100);"`                                // a short and optional name
	Priority          uint                   `gorm:"index; default:1;" json:"priority"`                             // the priority of the job (1 lowest priority, ... 10 high priority)
	Attempts          uint                   `gorm:"index; default:0;" json:"attempts"`                             // the attempts that have been tried in the past
	CurrentWorkerUID  *uuid.UUID             `gorm:"index; type:varchar(36); default:null" json:"currentWorkerUID"` // the uuid of the the worker bot that is currently working in this job
	WorkerInstance    string                 `gorm:"index" json:"workerInstance"`                                   // the type of worker that should work on this job
	Parameters        map[string]string      `gorm:"-" json:"parameters"`                                           // the search parameters the worker should use (e.g. company names ... )
	ParametersJSON    datatypes.JSON         `gorm:"type:jsonb;" json:"-"`                                          // json
	Tasks             []JobTask              `gorm:"-" json:"tasks"`                                                // the way the worker should interact (e.g. click on that, extract this, store it in this bucket, ...)
	TasksJSON         datatypes.JSON         `gorm:"type:jsonb;" json:"-"`                                          // json
	CurrentStatus     map[string]interface{} `gorm:"-" json:"currentStatus"`                                        // the current status of the job
	CurrentStatusJSON datatypes.JSON         `gorm:"type:jsonb;" json:"-"`                                          // json
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
