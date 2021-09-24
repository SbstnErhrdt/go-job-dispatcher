# Cleaning cron job

Cleans the database and resets stalled jobs.

```shell
kubectl apply -f cronjob.yml
```


```yml
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: CronJob
metadata:
  name: go-job-dispatcher-cleaner-job # Name of the Kubernetes resource
  namespace: go-job-dispatcher # Name of the Kubernetes namespace
  labels:
    app: go-job-dispatcher
spec:
  schedule: "*/5 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          imagePullSecrets:
            - name: go-registry
          containers:
            - name: go-jobs-clean-stalled-worker
              image: registry.url.com/job_dispatcher_cleaner:latest # TODO: change
              envFrom:
                - secretRef:
                    name: go-job-dispatcher-env-secret
          restartPolicy: OnFailure
EOF
```