package main

import (
	"workshops/workshop-04/pkg/hello"
)

var (
	version = "unknown"
)

func main() {
	hello.Run(version, "world")
}
