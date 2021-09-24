package redis_job_dispatcher

import (
	"github.com/SbstnErhrdt/env"
	"strings"
)

const (
	delimiter     = ":"
	keyTodo       = "todo"
	keyJobs       = "jobs"
	keyQueue      = "queue"
	keyDoing      = "doing"
	keyDone       = "done"
	keyWorkerJobs = "worker_jobs"
)

var (
	prefix = env.FallbackEnvVariable("JOB_DISPATCHER_REDIS_PREFIX", "job_dispatcher")
)

func init() {
	// add : to the end of the prefix
	// if there is not already a :
	if prefix[len(prefix)-1] != ':' {
		prefix = prefix + ":"
	}
}

// GetKeyPrefix returns the prefix of the redis keys
func GetKeyPrefix() string {
	return prefix
}

func GenerateKeyTodoMap(workerInstance string) string {
	return prefix + strings.ToLower(workerInstance) + delimiter + keyTodo
}

func GenerateKeyJobsMap() string {
	return prefix + keyJobs
}

func GenerateKeyTodoQueue(workerInstance string) string {
	return prefix + strings.ToLower(workerInstance) + delimiter + keyQueue
}

func GenerateKeyDoingMap(workerInstance string) string {
	return prefix + strings.ToLower(workerInstance) + delimiter + keyDoing
}

func GenerateKeyDoneMap(workerInstance string) string {
	return prefix + strings.ToLower(workerInstance) + delimiter + keyDone
}

func GenerateKeyWorkerJobs() string {
	return prefix + keyWorkerJobs
}
