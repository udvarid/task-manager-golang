package initJob

import (
	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/service"
)

func DoInitJob(config *model.Configuration) {
	overdueTasks, _ := service.GetAllTasks("")
	for _, task := range overdueTasks {
		communicator.SendMessage(config, task.Owner, task.Task)
	}
}
