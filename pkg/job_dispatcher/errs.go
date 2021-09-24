package job_dispatcher

import "errors"

var (
	NoNewJobs = errors.New("can not find a job")
)
