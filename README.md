
## Docker

Using docker-compose to build `wappsto-kafka-connector` with all its dependencies.
```
docker-compose up -d wappsto-connector
```
Command shold build images and spin all containers.

Check for running containers and bash into docker container
```
docker container ls
docker exec -it docker_wappsto-connector_1 bash
```


Follow appilication logs:
```
docker container logs docker_wappsto-connector_1 -f
```

