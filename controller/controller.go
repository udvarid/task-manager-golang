package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/udvarid/task-manager-golang/authenticator"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/service"
)

var (
	activeConfiguration = &model.Configuration{}
)

type GetSession struct {
	Id string `json:"id"`
}

func Init(config *model.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/validate/", validate)
	router.GET("/checkin/:id/:session", checkInTask)
	router.GET("/task/", taskPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.POST("/prolong/:id/:prolong_days", prolongTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

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

	isValidatedInTime := authenticator.Validate(activeConfiguration, getSession.Id, newSession)

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
		owner := validateSession(c)
		service.DeleteTask(deleteId, owner)
		redirectTo(c, "/task")
	}
}

func prolongTask(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		if prolong_days, err := strconv.Atoi(c.Param("prolong_days")); err == nil {
			owner := validateSession(c)
			service.ProlongTask(id, prolong_days, owner)
			redirectTo(c, "/task")
		}
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
