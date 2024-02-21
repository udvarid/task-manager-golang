package main

import (
	"flag"

	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/controller"
	"github.com/udvarid/task-manager-golang/repository"
)

var config = configuration.Configuration{}

func main() {
	configFile := flag.String("config", "conf.json", "the Json file contains the configurations")
	environment := flag.String("environment", "local", "where do we run tha application, local or on internet?")
	remoteAddress := flag.String("remote_address", "https://task-manager-golang.fly.dev/", "remote address of the application")
	dbName := flag.String("db_name", "my.db", "name of db name")

	flag.Parse()

	config = configuration.InitConfiguration(*configFile)
	config.Environment = *environment
	config.RemoteAddress = *remoteAddress
	config.DbName = *dbName
	repository.Init(&config)

	controller.Init(&config)
}
