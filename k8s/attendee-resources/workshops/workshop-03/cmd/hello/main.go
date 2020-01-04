package main

import (
	"workshops/workshop-03/pkg/hello"
)

var (
	version = "unknown"
)

func main() {
	hello.Run(version, "world")
}
