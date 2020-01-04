package main

import (
	"workshops/workshop-12/pkg/hello"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version = "unknown"
)

func main() {
	loadConfig()

	if viper.GetBool("oneshot-task") {
		hello.RunTask(viper.GetString("oneshot-task-endpoint"))
	} else {
		hello.RunServer(version, viper.GetString("name"))
	}
}

func loadConfig() {
	/* Defaults */
	viper.SetDefault("name", "world")

	/* Environment */
	viper.SetEnvPrefix("hello")
	viper.AutomaticEnv() /* e.g. HELLO_NAME */

	// Command-line args
	pflag.String("name", "world", "addressee") // i.e. --name
	pflag.Bool("oneshot-task", false, "run a task rather than the server")
	pflag.String("oneshot-task-endpoint", "https://webhook.site/xxx...", "endpoint to hit with the one-shot task")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
