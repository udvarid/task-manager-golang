package initJob

import (
	"time"

	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/service"
)

func DoInitJob(config *model.Configuration) {
	overdueTasks, _ := service.GetAllTasks("")
	for _, task := range overdueTasks {
		_, _, lastWarnedDay := task.LastWarnedTime.Date()
		_, _, todayDay := time.Now().Date()
		if lastWarnedDay != todayDay {
			service.UpdateWarningDate(task.ID)
			communicator.SendMessage(config, task.Owner, task.Task, task.DeadLStr)
		}
	}
}
