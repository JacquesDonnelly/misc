package main

import (
	"workshops/workshop-06/pkg/hello"
)

var (
	version = "unknown"
)

func main() {
	hello.Run(version, "world")
}
