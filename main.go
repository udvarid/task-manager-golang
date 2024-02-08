package main

import (
	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/controller"
)

func main() {
	controller.Init(configuration.InitConfiguration())
}
