# Job-Dispatcher

This repository contains the code of an `atomic like` jobs dispatcher. It is backed with sql database or with a redis
key-value store.

## Status

Work in progress

## Features

* Counts the attempts of jobs
* Heartbeat to store the current status of the job
* Priority
* Multi-Instance Jobs

## Dependencies

* SQL
* Redis

# Environment Variables

```
# SQL Database
SQL_TYPE=MYSQL
SQL_HOST=localhost
SQL_USER=root
SQL_PASSWORD=test
SQL_PORT=3306
SQL_DATABASE=test

# Clean jobs after n seconds / minutes ... SQL INTERVAL
CLEAN_STALLED_JOBS_INTERVAL=20 SECOND

# Application port
PORT=18989

# Job dispatcher prefix for redis
JOB_DISPATCHER_REDIS_PREFIX=job_dispatcher

# Redis Key Value Store
REDIS_HOST=192.168.157.33
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

```

# Tests

Run the tests via the script

```
sh tests.sh
```

or via the golang command

```
go test ./...
```

To test the redis version just add the following environment variable

```
export TEST_ENV=REDIS
```

# Client

Create a new job
```go
// create uuid
uid := uuid.New()
// init job dispatcher client
client := job_dispatcher.NewClient(
    https: //job-dispatcher.endpoint.net,
    "web-crawler",
    uid,
)
// create new job
newJob := job_dispatcher.NewJobDTO{
    Name:           "Search",
    Priority:       10,
    WorkerInstance: "web-crawler",
    Tasks: []job_dispatcher.JobTask{
        {
            Version: "0.1",
            Name:    "daily search",
            Type:    "search",
            Execute: map[string]Interface{}{
                "Name": "Company Name",
            },
        },
    },
}
// send job to dispatcher
res, err := client.CreateJob(newJob)
```


Receive the latest job in the queue
```go
// create uuid
uid := uuid.New()
// init job dispatcher client
client := job_dispatcher.NewClient(
    https: //job-dispatcher.endpoint.net,
    "web-crawler",
    uid,
)
// get the latest job from the queue
err := client.GetJob("job-type")
if err != nil {
    return
}

// start the job
err := client.StartCurrentJob()

// send a heartbeat with metrics of the job to the queue
err := client.HeartBeat(...)

// release the current job if something fails
_ = client.ReleaseCurrentJob()

// complete the jobs
_ := client.MarkCurrentJobAsCompleted()
```

# Deployment

Via docker
````shell
docker run -p 8080:8080 ese7en/go-job-dispatcher
````


* [Kubernetes](deployments/kubernetes)
* [Docker-Compose](deployments/kubernetes)