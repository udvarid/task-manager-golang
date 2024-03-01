package communicator

import (
	"log"
	"net/http"
	"net/smtp"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/udvarid/task-manager-golang/configuration"
)

var verifier = emailverifier.NewVerifier()

func SendMessageWithLink(config *configuration.Configuration, toAddress string, toLink string) {
	ret, err := verifier.Verify(toAddress)
	if err != nil || !ret.Syntax.Valid {
		sendNtfy(toAddress, "CheckInPls!", toLink)
	} else {
		msg := []byte("To: " + toAddress + "\r\n" +
			"Subject: Please check in!\r\n" +
			"\r\n" +
			"Here is the link\r\n" +
			toLink)
		sendMail(config, toAddress, msg)
	}
}

func SendMessage(config *configuration.Configuration, toAddress string, task string) {
	ret, err := verifier.Verify(toAddress)
	if err != nil || !ret.Syntax.Valid {
		sendNtfyMessage(toAddress, "Overdue: "+task)
	} else {
		msg := []byte("To: " + toAddress + "\r\n" +
			"Subject: Overdue task!\r\n" +
			"\r\n" +
			task)
		sendMail(config, toAddress, msg)
	}
}

func sendMail(config *configuration.Configuration, toAddress string, message []byte) {
	auth := smtp.PlainAuth("", config.Mail_from, config.Mail_psw, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, config.Mail_from, []string{toAddress}, message)
	if err != nil {
		log.Print(err)
	}
}

func sendNtfy(channel string, msg string, address string) {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/"+channel, strings.NewReader(msg))
	req.Header.Set("Actions", "http, Confirm!, "+address+", method=GET")
	http.DefaultClient.Do(req)
}

func sendNtfyMessage(channel string, msg string) {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/"+channel, strings.NewReader(msg))
	req.Header.Set("Title", "Overdue task")
	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "warning,skull")
	http.DefaultClient.Do(req)
}
