package main

import (
	"embed"
	"encoding/json"
	"flag"

	"github.com/udvarid/task-manager-golang/controller"
	"github.com/udvarid/task-manager-golang/initJob"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository/taskRepository"
)

var config = model.Configuration{}

//go:embed resources
var f embed.FS

func main() {
	configFile := flag.String("config", "conf.json", "the Json file contains the configurations")
	environment := flag.String("environment", "local", "where do we run tha application, local or on internet?")
	remoteAddress := flag.String("remote_address", "https://task-manager-golang.fly.dev/", "remote address of the application")
	flag.Parse()

	configFileInString, _ := f.ReadFile("resources/" + *configFile)
	json.Unmarshal([]byte(configFileInString), &config)

	config.Environment = *environment
	config.RemoteAddress = *remoteAddress
	taskRepository.Init()

	initJob.DoInitJob(&config)

	controller.Init(&config)
}
