package tests

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// test the health endpoint
func TestHealth(t *testing.T) {
	ass := assert.New(t)
	resp, err := http.Get(Server.URL + "/health")
	ass.NoError(err)
	ass.Equal(200, resp.StatusCode)
	return
}

// test the liveness endpoint
func TestLiveness(t *testing.T) {
	ass := assert.New(t)
	resp, err := http.Get(Server.URL + "/liveness")
	ass.NoError(err)
	ass.Equal(200, resp.StatusCode)

	// Kill the database connection
	connections.SQLClient = nil
	resp, err = http.Get(Server.URL + "/liveness")
	ass.NoError(err)
	ass.NotEqual(200, resp.StatusCode)

	// Connect again
	connections.ConnectToSQL()
	resp, err = http.Get(Server.URL + "/liveness")
	ass.NoError(err)
	ass.Equal(200, resp.StatusCode)
	return
}
