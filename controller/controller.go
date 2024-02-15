package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/udvarid/task-manager-golang/authenticator"
	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/service"
)

var activeConfiguration = &configuration.Configuration{}

type NewTask struct {
	Task string `json:"task"`
}

type GetSession struct {
	Id string `json:"id"`
}

func Init(config *configuration.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/validate/", validate)
	router.GET("/task/", taskPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

// TODO
// 1, authentikáció megoldása (ez valószínűleg localon nem is tesztelhető, webhook kell hozzá)
// 2, Db implementáció
// 3, fly.io.n tesztelni

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

	if activeConfiguration.Environment == "local" {
		fmt.Println("Local environment, validation prcess skipped")
	} else {
		//  ezek majd az authentikációs service-be menjenek át
		//  communicator.SendNtfy("donat1977", "hello-bello", "http://localhost:8080/")
		//	communicator.SendMail(activeConfiguration, "udvarid@hotmail.com", []byte("Hello"))
	}

	newSession := authenticator.GiveSession(getSession.Id)

	c.SetCookie("id", getSession.Id, 3600, "/", "localhost", false, true)
	c.SetCookie("session", newSession, 3600, "/", "localhost", false, true)
	redirectTo(c, "/task")
}

func taskPage(c *gin.Context) {
	id := validateSession(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
		"tasks": service.GetAllTasks(id),
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
	var newTask NewTask
	c.BindJSON(&newTask)
	id := validateSession(c)
	service.AddTask(newTask.Task, id)
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
