package main

import (
	"fmt"
	"workshops/workshop-17/pkg/hello"

	"github.com/fsnotify/fsnotify"
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
		hello.RunServer(version, viper.GetString("name"), viper.GetString("redis_url"))
	}
}

func loadConfig() {
	/* Least to most priority */

	/* Defaults */
	viper.SetDefault("name", "world")
	viper.SetDefault("redis_url", "localhost:6379")

	/* Config file */
	viper.SetConfigName("config") /* e.g. config.yaml, config.json */
	viper.AddConfigPath("$HOME/.hello/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// Also goes off for no config file
		// fmt.Fatalf("Fatal error in config file: %s \n", err)
	}

	/* ...with auto-reload */
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

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
