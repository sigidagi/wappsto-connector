## Intro

### Start Wappsto connector in docker

```
docker-compose up
```



### A UI for Data Streaming

Redpanda console will then be run as it would be executed on the host machine.
```
docker run --network=host -p 8080:8080 -e KAFKA_BROKERS=localhost:29092 docker.redpanda.com/vectorized/console:master-0a8fce8
```


### Docker command helpers

Remove old and dangling images:
```
docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
```

Using docker-compose to build `wappsto-kafka-connector` with all its dependencies.
```
docker-compose up -d wappsto-connector
```

Check for running containers and bash into docker container
```
docker container ls
docker exec -it docker_wappsto-connector_1 bash
```

Follow appilication logs:
```
docker container logs docker_wappsto-connector_1 -f
```

Stop all container and rebuild again:
```
docker stop $(docker ps -q)
docker-compose up --build
```
