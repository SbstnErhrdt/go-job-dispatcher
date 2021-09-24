# Redis

* Key/Value Store: where all the jobs are stored
* Todo Queue: Contains the ids of the jobs that are currently todo

## Data Structure

### Jobs

| Key         |      Val             |  
|-------------|:--------------------:|
| uuid of job |  json object of jobs |

### {worker_instance} : todo

A hash map with the todo jobs

| Key         |      Val        |  
|-------------|:---------------:|
| uuid of job |  1              |

### {worker_instance} : queue

A simple queue with the uids of the jobs

| Queue Items |
|-------------|
| uuid of job |

### {worker_instance} : doing

A hash map with the current jobs

| Key         |      Val        |  
|-------------|:---------------:|
| uuid of job |  unix-timestamp |

### {worker_instance} : done

A hash map with the done jobs

| Key         |      Val        |  
|-------------|:---------------:|
| uuid of job |  1              |