package communicator

import (
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/udvarid/task-manager-golang/configuration"
)

func SendMail(config *configuration.Configuration, toAddress string, message []byte) {
	auth := smtp.PlainAuth("", config.Mail_from, config.Mail_psw, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, config.Mail_from, []string{toAddress}, message)
	if err != nil {
		log.Print(err)
	}
}

func SendNtfy(channel string, msg string, address string) {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/"+channel, strings.NewReader(msg))
	req.Header.Set("Actions", "http, Confirm!, "+address+", method=GET")
	http.DefaultClient.Do(req)
}

func SendNtfyMessage(channel string, msg string) {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/"+channel, strings.NewReader(msg))
	req.Header.Set("Title", "Overdue task")
	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "warning,skull")
	http.DefaultClient.Do(req)
}
