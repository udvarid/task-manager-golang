package main

import (
	"flag"

	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/controller"
)

var config = configuration.Configuration{}

func main() {
	configFile := flag.String("config", "conf.json", "the Json file contains the configurations")
	environment := flag.String("environment", "local", "where do we run tha application, local or on internet?")
	flag.Parse()

	config = configuration.InitConfiguration(*configFile)
	config.Environment = *environment
	controller.Init(&config)
}
