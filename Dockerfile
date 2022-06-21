FROM golang:1.18-alpine AS development

ENV PROJECT_PATH=/wappsto-kafka-connector
ENV PATH=$PATH:$PROJECT_PATH/bin
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache ca-certificates tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make

FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata bash
COPY --from=development /wappsto-kafka-connector/bin/wappsto-kafka-connector /usr/bin/wappsto-kafka-connector
COPY --from=development /wappsto-kafka-connector/wappsto-kafka-connector.toml /etc/wappsto-kafka-connector/wappsto-kafka-connector.toml

ENTRYPOINT ["/usr/bin/wappsto-kafka-connector"]
