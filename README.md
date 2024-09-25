# Airport API
# Containerized Go App
```shell
docker build -t sre_airports_api .
docker run -p 8080:8080 sre_airports_api
```


# Terraform Cloud Storage Update
Check all the steps.
https://github.com/zaxxio/sre_airports_api/actions/runs/11040296438/job/30667945334#step:1:1


# Docker Hub Latest Build only one. Zero versioning
https://hub.docker.com/repository/docker/polymerpro/sre-airports-api

# Kubernetes Endpoint to check the services
http://34.135.84.233:8080/airports


# However i've not added loadbalancer for traffic spilting and kubernetes automatic deployment.
# I'm new GCP and but has experience in AWS.
