package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/udvarid/task-manager-golang/authenticator"
	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/service"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	verifier            = emailverifier.NewVerifier()
	activeConfiguration = &configuration.Configuration{}
)

type GetSession struct {
	Id string `json:"id"`
}

func Init(config *configuration.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/validate/", validate)
	router.GET("/checkin/:id/:session", checkInTask)
	router.GET("/task/", taskPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

// TODO
// 3, a scheduled futtatás kapcsán a main-ből induljon el egy task, ami ellenőrzi az adatbázist, hogy van e lejáró
// 4. Lehessen taskot hosszabbítani 1 nap/1 héttel/1 hónappal (+1-1 gomb)
// 5, Go embed feature-ét használni, a templatek és a conf.json file-ra
// 6. Refactor: Belépéskori validálást áthelyezni az authentikátorba
// 7, Kicsinosítani a frontendet
// 8, Fly-on tesztelni a mount-ot és a scheduled futtatás

func startPage(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", "localhost", false, true)
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.HTML(http.StatusOK, "start.html", gin.H{
		"title": "Home Page",
	})
}

func validate(c *gin.Context) {
	var getSession GetSession
	c.BindJSON(&getSession)
	newSession, err := authenticator.GiveSession(getSession.Id)
	if err != nil {
		redirectTo(c, "/")
	}

	isValidatedInTime := false
	if activeConfiguration.Environment == "local" {
		fmt.Println("Local environment, validation process skipped")
		isValidatedInTime = true
	} else {
		linkToSend := activeConfiguration.RemoteAddress + "checkin/" + getSession.Id + "/" + newSession
		ret, err := verifier.Verify(getSession.Id)
		if err != nil || !ret.Syntax.Valid {
			communicator.SendNtfy(getSession.Id, "CheckInPls!", linkToSend)
		} else {
			msg := []byte("To: " + getSession.Id + "\r\n" +
				"Subject: Please check in!\r\n" +
				"\r\n" +
				"Here is the link\r\n" +
				linkToSend)
			communicator.SendMail(activeConfiguration, getSession.Id, msg)
		}

		foundChecked := make(chan string)
		timer := time.NewTimer(60 * time.Second)
		go func() {
			for {
				time.Sleep(1 * time.Second)
				isCheckedAlready := authenticator.IsChecked(getSession.Id, newSession)
				if isCheckedAlready {
					foundChecked <- "one"
				}

			}
		}()
		select {
		case <-foundChecked:
			fmt.Println("Id is validated")
			isValidatedInTime = true
		case <-timer.C:
			fmt.Println("Id is not validated in time")
		}
	}
	if isValidatedInTime {
		c.SetCookie("id", getSession.Id, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		c.SetCookie("session", newSession, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		redirectTo(c, "/task")
	} else {
		redirectTo(c, "/")
	}

}

func checkInTask(c *gin.Context) {
	authenticator.CheckIn(c.Param("id"), c.Param("session"))
}

func taskPage(c *gin.Context) {
	id := validateSession(c)
	overDue, normal := service.GetAllTasks(id)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":  "Home Page",
		"tasksO": overDue,
		"tasksN": normal,
	})
}

func deleteTask(c *gin.Context) {
	if deleteId, err := strconv.Atoi(c.Param("delete_id")); err == nil {
		validateSession(c)
		service.DeleteTask(int(deleteId))
		redirectTo(c, "/task")
	}
}

func addTask(c *gin.Context) {
	var newTask model.NewTask
	c.BindJSON(&newTask)
	id := validateSession(c)
	service.AddTask(newTask, id)
	redirectTo(c, "/task")
}

func newTask(c *gin.Context) {
	validateSession(c)
	c.HTML(http.StatusOK, "addNew.html", gin.H{})
}

func redirectTo(c *gin.Context, path string) {
	location := url.URL{Path: path}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func validateSession(c *gin.Context) string {
	id_cookie, err := c.Cookie("id")
	isMissingCookie := false
	if err != nil {
		isMissingCookie = true
	}
	session_cookie, err := c.Cookie("session")
	if err != nil {
		isMissingCookie = true
	}
	if isMissingCookie {
		redirectTo(c, "/")
	}
	isValid := authenticator.IsValid(id_cookie, session_cookie)
	if !isValid {
		redirectTo(c, "/")
	}
	return id_cookie
}
