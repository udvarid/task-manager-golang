package initJob

import (
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/service"
)

var verifier = emailverifier.NewVerifier()

func DoInitJob(config *configuration.Configuration) {
	overdueTasks, _ := service.GetAllTasks("")
	for _, task := range overdueTasks {
		ret, err := verifier.Verify(task.Owner)
		if err != nil || !ret.Syntax.Valid {
			communicator.SendNtfyMessage(task.Owner, "Overdue: "+task.Task)
		} else {
			msg := []byte("To: " + task.Owner + "\r\n" +
				"Subject: Overdue task!\r\n" +
				"\r\n" +
				task.Task)
			communicator.SendMail(config, task.Owner, msg)
		}

	}
}
