BINARY_NAME=wappsto-kafka-connector
VERSION := 0.1.0

all: build

build:
	@echo "Compiling source"
	mkdir -p bin
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o bin/${BINARY_NAME} cmd/wappsto-kafka-connector/main.go

run: build
	@echo "Starting Wappsto connector to Kafka service"
	./bin/${BINARY_NAME}

clean:
	@echo "cleaning builds"
	go clean
	rm ./bin/${BINARY_NAME}
