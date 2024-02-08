package controller

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"net/smtp"

	"github.com/gin-gonic/gin"

	"github.com/udvarid/task-manager-golang/configuration"
	"github.com/udvarid/task-manager-golang/service"
)

var activeConfiguration = configuration.Configuration{}

func Init(config configuration.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

func startPage(c *gin.Context) {
	myMail := "donat.udvari@gmail.com"
	auth := smtp.PlainAuth("", myMail, activeConfiguration.Mail_psw, "smtp.gmail.com")
	to := []string{"udvarid@hotmail.com"}
	msg := []byte("To: kate.doe@example.com\r\n" +
		"Subject: Why aren’t you using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, myMail, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
		"tasks": service.GetAllTasks(),
	})
}

func deleteTask(c *gin.Context) {
	if deleteId, err := strconv.Atoi(c.Param("delete_id")); err == nil {
		service.DeleteTask(int(deleteId))
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func addTask(c *gin.Context) {
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func newTask(c *gin.Context) {
	c.HTML(http.StatusOK, "addNew.html", gin.H{})
}
