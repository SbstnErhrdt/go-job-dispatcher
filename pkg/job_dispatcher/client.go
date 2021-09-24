package job_dispatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Client struct {
	Endpoint       string
	UUID           uuid.UUID
	WorkerInstance string
	CurrentJob     *Job
	Bulk           bool
}

// NewClient inits a new client
func NewClient(endpoint string, workerInstance string, uuid uuid.UUID) Client {
	// fix endpoint trailing slash
	l := len(endpoint)
	if endpoint[l-1] == '/' {
		endpoint = endpoint[:l-1]
	}
	// init client
	return Client{
		Endpoint:       endpoint,
		UUID:           uuid,
		WorkerInstance: workerInstance,
		Bulk:           false,
	}
}

// UseBulk uses the bulk endpoint
func (c *Client) UseBulk() *Client {
	c.Bulk = true
	return c
}

// UseDefault uses the default endpoint
func (c *Client) UseDefault() *Client {
	c.Bulk = false
	return c
}

func (c *Client) GetEndpoint() string {
	if c.Bulk {
		return c.Endpoint + "/bulk-jobs"
	} else {
		return c.Endpoint + "/jobs"
	}
}

// GetJob gets the latest job with the highest priority
func (c *Client) GetJob(additionalInstances []string) (err error) {
	// create payload
	payload := map[string]interface{}{}
	// add additional instances
	if len(additionalInstances) > 0 {
		payload["instances"] = additionalInstances
	}
	// marshal payload
	jsonPayload, _ := json.Marshal(payload)
	// build url with uuid of the worker
	url := c.GetEndpoint() + "/get/" + c.WorkerInstance + "/" + c.UUID.String()
	// send request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// initialize http client
	client := &http.Client{}
	// execute the request and provide the response
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// if there is no job
	if resp.StatusCode == 204 {
		c.CurrentJob = nil
		err = nil
		return
	}
	if resp.StatusCode > 399 {
		err = errors.New("no 20X status")
		return
	}
	// parse the response
	var jobResponse ResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&jobResponse)
	if err != nil {
		return
	}
	// assign the job to the client
	log.Println("job assigned", jobResponse.Job.Name, jobResponse.Job.UUID.String())
	c.CurrentJob = jobResponse.Job
	return
}

// CreateJob creates a single job
func (c *Client) CreateJob(newJob NewJobDTO) (job *Job, err error) {
	// marshal payload
	jsonPayload, _ := json.Marshal(newJob)
	// send request
	resp, err := http.Post(c.GetEndpoint(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	if resp.StatusCode > 399 {
		err = errors.New("no 20X status")
		return
	}
	response := ResponseDTO{}
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	return response.Job, nil
}

// CreateJobs creates multiple job at once
func (c *Client) CreateJobs(newJobs []NewJobDTO) (results []*Job, err error) {
	// marshal payload
	jsonPayload, _ := json.Marshal(newJobs)
	// send request
	resp, err := http.Post(c.GetEndpoint()+"/bulk", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	if resp.StatusCode > 399 {
		err = errors.New("no 20X status")
		return
	}
	response := ResponseMultipleDTO{}
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	return response.Jobs, nil
}

// sendReq is a helper function that sends http put req to different endpoints
func (c *Client) sendReq(reqType string) (err error) {
	if c.CurrentJob == nil {
		return errors.New("no current job set")
	}
	// marshal payload
	// build url with uuid of the worker
	url := c.GetEndpoint() + "/" + reqType + "/" + c.CurrentJob.UUID.String()
	// send request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(nil))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// initialize http client
	client := &http.Client{}
	// execute the request and provide the response
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode > 399 {
		err = errors.New("no 20X status")
		return
	}
	response := ResponseDTO{}
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	// assign job
	c.CurrentJob = response.Job
	return
}

// HeartBeat sends the heartbeat to the backend
func (c *Client) HeartBeat(status map[string]interface{}) (err error) {
	if c.CurrentJob == nil {
		return errors.New("no current job set")
	}
	// marshal payload
	// build url with uuid of the worker
	url := c.GetEndpoint() + "/heartbeat/" + c.CurrentJob.UUID.String()
	// send request
	payload, err := json.Marshal(status)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// initialize http client
	client := &http.Client{}
	// execute the request and provide the response
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode > 399 {
		err = errors.New("no 20X status")
		return
	}
	response := ResponseDTO{}
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	// assign job
	c.CurrentJob = response.Job
	return
}

// StartCurrentJob marks the current job as started
func (c *Client) StartCurrentJob() (err error) {
	return c.sendReq("start")
}

// MarkCurrentJobAsCompleted marks the current job as completed
func (c *Client) MarkCurrentJobAsCompleted() (err error) {
	err = c.sendReq("complete")
	if err != nil {
		return err
	}
	// set current job to nil
	c.CurrentJob = nil
	return
}

// ReleaseCurrentJob releases the current job
func (c *Client) ReleaseCurrentJob() (err error) {
	err = c.sendReq("release")
	if err != nil {
		return err
	}
	// set current job to nil
	c.CurrentJob = nil
	return
}
