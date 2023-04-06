package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func GetStalledJobs() (stalledJobs []uuid.UUID, err error) {
	instances, err := GetInstances()
	if err != nil {
		log.WithError(err).Error("Could not get instances")
		return
	}
	// iterate over all instances
	for _, instance := range instances {
		// Get all jobs in doing and check if they are stalled
		resCMD := connections.RedisClient.HGetAll(context.TODO(), GenerateKeyDoingMap(instance))
		warn := resCMD.Err()
		if warn != nil {
			log.WithField("instance", instance).Warn(warn)
			continue
		}
		res, warn := resCMD.Result()
		if warn != nil {
			log.WithField("instance", instance).Warn(warn)
			continue
		}
		for key, val := range res {
			// get the time from the unix string
			i, warnUnix := strconv.ParseInt(val, 10, 64)
			if warnUnix != nil {
				log.WithField("instance", instance).Warn(warnUnix)
				continue
			}
			lastInteraction := time.Unix(i, 0).UTC()
			// compare now and last interaction
			now := time.Now().UTC()
			if now.Sub(lastInteraction).Minutes() > 10 {
				// if there has been a stalled job
				// parse the uuid of the stalled job
				jobUID, parseWarn := uuid.Parse(key)
				if parseWarn != nil {
					log.WithField("instance", instance).Warn(parseWarn)
					continue
				}
				stalledJobs = append(stalledJobs, jobUID)
			}
		}
	}
	return
}
