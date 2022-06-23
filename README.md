## Intro

### Start Wappsto connector in docker

```
docker-compose up
```

### A UI for Data Streaming

[Redpanda Console](https://github.com/redpanda-data/console)

Redpanda console will .
```
docker run --network=host -p 8080:8080 -e KAFKA_BROKERS=localhost:29092 docker.redpanda.com/vectorized/console:master-0a8fce8
```
