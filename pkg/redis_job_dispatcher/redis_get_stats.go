package redis_job_dispatcher

import (
	"context"
	"github.com/SbstnErhrdt/go-job-dispatcher/connections"
	"github.com/SbstnErhrdt/go-job-dispatcher/pkg/job_dispatcher"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strings"
)

// GetKeys returns all the key from the database
func GetKeys() (results []string, err error) {
	resCMD := connections.RedisClient.Keys(context.TODO(), "*")
	if resCMD.Err() != nil {
		err = resCMD.Err()
		log.WithError(err).Error("Could not get keys from redis")
		return
	}
	results, err = resCMD.Result()
	if err != nil {
		log.WithError(err).Error("Could not get keys from redis")
		return
	}
	return
}

// GetStats calculates the stats by iterating over the keys in the redis store
func GetStats() (results []job_dispatcher.Stats, err error) {
	stats := map[string]*job_dispatcher.Stats{}
	keys, err := GetKeys()
	if err != nil {
		log.WithError(err).Error("Could not get keys from redis")
		return
	}
	// start tx
	tx := connections.RedisClient.TxPipeline()
	// iterate over all the keys
	for _, key := range keys {
		// get the key fragments
		keyFragments := strings.Split(key, ":")
		l := len(keyFragments)
		if l > 2 {
			lastFragment := keyFragments[l-1]
			// add instance to struct
			instance := keyFragments[l-2]
			if _, ok := stats[instance]; !ok {
				// if there is not already the instance in the map add it
				s := job_dispatcher.Stats{
					WorkerInstance: instance,
				}
				stats[instance] = &s
			}
			switch lastFragment {
			case "todo":
				tx.HLen(context.TODO(), key)
			case "doing":
				tx.HLen(context.TODO(), key)
			case "done":
				tx.HLen(context.TODO(), key)
			case "queue":
				tx.LLen(context.TODO(), key)
			}
		}
	}
	// execute tx
	resCMD, err := tx.Exec(context.TODO())
	if err != nil {
		log.WithError(err).Error("Could not get stats from redis")
		return
	}

	for _, res := range resCMD {

		key := res.Args()[1].(string)

		keyFragments := strings.Split(key, ":")
		l := len(keyFragments)
		lastFragment := keyFragments[l-1]
		// add instance to struct
		instance := keyFragments[l-2]
		switch lastFragment {
		case "todo":
			stats[instance].Todo = int(res.(*redis.IntCmd).Val())
		case "doing":
			stats[instance].Active = int(res.(*redis.IntCmd).Val())
		case "done":
			stats[instance].Done = int(res.(*redis.IntCmd).Val())
		case "queue":
			stats[instance].Todo = int(res.(*redis.IntCmd).Val())
		}
	}

	// transform map to array
	for _, v := range stats {
		v.Total = v.Todo + v.Active + v.Done
		results = append(results, *v)
	}
	return
}

// GetInstances returns the instances present in the database
func GetInstances() (result []string, err error) {
	res := map[string]struct{}{}
	// get all keys from the database
	keys, err := GetKeys()
	if err != nil {
		return
	}
	// iterate over all the keys
	for _, key := range keys {
		// get the key fragments
		keyFragments := strings.Split(key, ":")
		l := len(keyFragments)
		if l > 2 {
			// add instance to struct
			instance := keyFragments[l-2]
			if _, ok := res[instance]; !ok {
				// if there is not already the instance in the map add it
				res[instance] = struct{}{}
			}
		}
	}
	// generate the array from the map
	result = []string{}
	for k := range res {
		result = append(result, k)
	}
	return
}
