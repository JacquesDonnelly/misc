package main

import (
	"fmt"
	"workshops/workshop-08/pkg/hello"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version = "unknown"
)

func main() {
	loadConfig()
	name := viper.GetString("name")

	hello.Run(version, name)
}

func loadConfig() {
	/* Least to most priority */

	/* Defaults */
	viper.SetDefault("name", "world")

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

	/* Command-line args */
	pflag.String("name", "world", "addressee") /* i.e. --name */
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
