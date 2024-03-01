package initJob

import (
	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/service"
)

func DoInitJob(config *configuration.Configuration) {
	overdueTasks, _ := service.GetAllTasks("")
	for _, task := range overdueTasks {
		communicator.SendMessage(config, task.Owner, task.Task)
	}
}
