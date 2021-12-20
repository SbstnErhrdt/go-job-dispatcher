package api

import (
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/redis_job_dispatcher"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/sql_job_dispatcher"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitServer() {
	log.Println("try to init server")
	r := InitServerEngine()
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("server successfully initialize")
}

func InitBackgroundJobs() {
	sql := sql_job_dispatcher.SqlService{}
	redis := redis_job_dispatcher.RedisService{}
	for {
		log.Info("clean up expired jobs")
		errSQL := sql.Clean()
		if errSQL != nil {
			log.Fatal("can not clean sql ", errSQL)
		}
		errRedis := redis.Clean()
		if errRedis != nil {
			log.Fatal("can not clean redis ", errRedis)
		}
		time.Sleep(time.Minute * 5)
	}
}

func InitServerEngine() (r *gin.Engine) {
	log.Println("try to init router")
	// init router
	r = gin.Default()
	// add the cors middleware
	r.Use(CORSMiddleware())

	// Kubernetes
	r.GET("/readiness", LivenessHandler)
	r.GET("/health", LivenessHandler)
	r.GET("/liveness", LivenessHandler)

	// Jobs

	jobs := r.Group("/jobs")
	// add the sql service
	jobs.Use(func(c *gin.Context) {
		c.Set(JobServiceKey, &sql_job_dispatcher.SqlService{})
		c.Next()
	})
	{
		jobs.POST("/", CreateJobHandler)
		jobs.POST("/bulk", BulkCreateJobsHandler)
		jobs.GET("/clean", CleanJobsHandler)
		jobs.PUT("/start/:uuid", StartJobHandler)
		jobs.PUT("/heartbeat/:uuid", HeartBeatJobHandler)
		jobs.PUT("/release/:uuid", ReleaseJobHandler)
		jobs.PUT("/complete/:uuid", CompleteJobHandler)
		jobs.PUT("/get/:instance/:uuid", GetJobsHandler)
		// Stats
		jobs.GET("/stats", ReadStatsJobHandler)
	}
	// Bulk jobs with redis
	bulkJobs := r.Group("/bulk-jobs")
	// add the sql service
	bulkJobs.Use(func(c *gin.Context) {
		c.Set(JobServiceKey, &redis_job_dispatcher.RedisService{})
		c.Next()
	})
	{
		bulkJobs.POST("/", CreateJobHandler)
		bulkJobs.POST("/bulk", BulkCreateJobsHandler)
		bulkJobs.GET("/clean", CleanJobsHandler)
		bulkJobs.PUT("/start/:uuid", StartJobHandler)
		bulkJobs.PUT("/heartbeat/:uuid", HeartBeatJobHandler)
		bulkJobs.PUT("/release/:uuid", ReleaseJobHandler)
		bulkJobs.PUT("/complete/:uuid", CompleteJobHandler)
		bulkJobs.PUT("/get/:instance/:uuid", GetJobsHandler)
		// Stats
		bulkJobs.GET("/stats", ReadStatsJobHandler)
	}
	log.Println("router successfully initialize")
	return
}
