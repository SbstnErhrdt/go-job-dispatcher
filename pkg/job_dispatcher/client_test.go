package job_dispatcher

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	ass := assert.New(t)
	c := NewClient("test", "test-instance", uuid.New())
	ass.Equal(c.GetEndpoint(), "test/jobs")
}

func TestClient_UseBulk(t *testing.T) {
	ass := assert.New(t)
	c := NewClient("test", "test-instance", uuid.New())
	c.UseBulk()
	ass.Equal(c.GetEndpoint(), "test/bulk-jobs")
	c.UseBulk()
	ass.Equal(c.GetEndpoint(), "test/bulk-jobs")
}

func TestClient_UseDefault(t *testing.T) {
	ass := assert.New(t)
	c := NewClient("test", "test-instance", uuid.New())
	c.UseBulk()
	ass.Equal(c.GetEndpoint(), "test/bulk-jobs")
	c.UseDefault()
	ass.Equal(c.GetEndpoint(), "test/jobs")
	c.UseDefault()
	ass.Equal(c.GetEndpoint(), "test/jobs")
}
