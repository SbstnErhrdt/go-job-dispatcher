package job_dispatcher

import "errors"

var (
	// ErrNoNewJobs is the error returned when there are no new jobs
	ErrNoNewJobs = errors.New("can not find a job")
)
