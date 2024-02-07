package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/udvarid/task-manager-golang/service"
)

func Init() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

func startPage(c *gin.Context) {
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
