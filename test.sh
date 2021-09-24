#!/bin/sh
# load env file
echo "Load env file"
if [ -f .env ]
then
  export $(cat .env | sed 's/#.*//g' | xargs)
fi
# test with sql
docker-compose up -d
echo "Test SQL"
go clean -testcache
go test ./tests
# test with redis
export TEST_ENV=REDIS
echo "Test Redis"
go clean -testcache
go test ./tests
go test ./pkg/redis_job_dispatcher
docker-compose down