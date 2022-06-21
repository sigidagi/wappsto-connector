package main

import (
	"wappsto-kafka-connector/cmd/wappsto-kafka-connector/cmd"
)

// TODO init all log levels here.

var version string // set by compiler

func main() {
	cmd.Execute(version)
}
