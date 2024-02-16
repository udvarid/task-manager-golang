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
	remoteAddress := flag.String("remote_address", "https://task-manager-golang.fly.dev/", "remote address of the application")

	flag.Parse()

	config = configuration.InitConfiguration(*configFile)
	config.Environment = *environment
	config.RemoteAddress = *remoteAddress
	controller.Init(&config)
}
