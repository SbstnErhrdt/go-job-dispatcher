package sql_job_dispatcher

import (
	"encoding/json"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Job extends the model so that it works with sql
type Job struct {
	job_dispatcher.Job
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
	ParametersJSON    datatypes.JSON `json:"-" gorm:"type:jsonb;"` // json
	TasksJSON         datatypes.JSON `json:"-" gorm:"type:jsonb;"` // json
	CurrentStatusJSON datatypes.JSON `json:"-" gorm:"type:jsonb;"` // json
}

func (wj *Job) BeforeSave(*gorm.DB) (err error) {
	// Marshal the structs into the json fields
	// parameters
	env, err := json.Marshal(wj.Parameters)
	if err != nil {
		return
	}
	err = wj.ParametersJSON.Scan(env)
	// tasks
	tasks, err := json.Marshal(wj.Tasks)
	if err != nil {
		return
	}
	err = wj.TasksJSON.Scan(tasks)
	// status
	status, err := json.Marshal(wj.CurrentStatus)
	if err != nil {
		return
	}
	err = wj.CurrentStatusJSON.Scan(status)
	return
}

func (wj *Job) AfterFind(*gorm.DB) (err error) {
	// Unmarshal the json fields back to the structs
	// parameters
	err = json.Unmarshal(wj.ParametersJSON, &wj.Parameters)
	if err != nil {
		return
	}
	// tasks
	err = json.Unmarshal(wj.TasksJSON, &wj.Tasks)
	if err != nil {
		return
	}
	// status
	err = json.Unmarshal(wj.CurrentStatusJSON, &wj.CurrentStatus)
	if err != nil {
		return
	}
	return
}
