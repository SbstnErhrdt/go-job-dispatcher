package api

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/gin-gonic/gin"
	"time"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
	return
}

func LivenessHandler(c *gin.Context) {
	type dbTest struct {
		Now time.Time
	}
	var res dbTest
	err := connections.SQLClient.Raw(`SELECT current_timestamp as now`).Scan(&res).Error
	if err != nil || res.Now.IsZero() {
		c.JSON(503, gin.H{
			"status": "DOWN",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "UP",
	})
	return
}
