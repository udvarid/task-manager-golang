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

var user = "donat1977"

type NewTask struct {
	Task  string `json:"task"`
	Owner string `json:"owner"`
}

type GetSession struct {
	Id string `json:"id"`
}

type ResultSession struct {
	Id      string `json:"id"`
	Session string `json:"session"`
}

func Init(config *configuration.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.POST("/validate", validate)
	router.GET("/task", taskPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

// TODO
// 1, Kell egy kiinduló lap, ahol ntfy azonosítót vagy emailt adunk meg
//    Azonosító megadása után egy random string generálódik és mentődik el egy map-ban (később db-ben)
//    Ha nem local-ban vagyunk, akkor ntfy vagy email-el validálni kell, ezt később, ehhez kell valószínűleg a webhook
//    Amúgy ez jó esetben visszakapja header-ben a user és a token-t
// 3, A task-os lista "/task" végponton legyen elérhető. Ide átadjuk a usernevet és a tokent is
// 4, Delete-nél headerbe jöjjön vissza a usernév és a token, hogy tudjuk validálni
// 5, Még nem tudom, hogy delete és new task után redirect-nél hogyan biztosítsuk, hogy user/token öröklést
// 6, Kell egy logout oldal
// 7, Db implementáció

func startPage(c *gin.Context) {
	c.HTML(http.StatusOK, "start.html", gin.H{
		"title": "Home Page",
	})
}

func validate(c *gin.Context) {
	var getSession GetSession
	c.BindJSON(&getSession)
	newSession := authenticator.GiveSession(getSession.Id)

	resultSession := ResultSession{Id: getSession.Id, Session: newSession}

	c.JSON(http.StatusOK, gin.H{"session": resultSession})
}

func taskPage(c *gin.Context) {
	//  ezek majd az authentikációs service-be menjenek át
	//  communicator.SendNtfy("donat1977", "hello-bello", "http://localhost:8080/")
	//	communicator.SendMail(activeConfiguration, "udvarid@hotmail.com", []byte("Hello"))

	// a megkapott user és token alapján authentikálni kell a user-t. Ha ok, akkor kiszoljáljuk
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
		"tasks": service.GetAllTasks(user),
	})
}

func deleteTask(c *gin.Context) {
	if deleteId, err := strconv.Atoi(c.Param("delete_id")); err == nil {
		// a megkapott user és token alapján authentikálni kell a user-t. Ha ok, akkor kiszoljáljuk
		service.DeleteTask(int(deleteId))
		location := url.URL{Path: "/task"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func addTask(c *gin.Context) {
	// itt kiszedjük az infókat ill. a username és token alapján authentikáljuk. ha okés, akkor kiszolgáljuk
	var newTask NewTask
	c.BindJSON(&newTask)
	service.AddTask(newTask.Task, newTask.Owner)
	mami := c.GetHeader("mami")
	fmt.Println(mami)
	location := url.URL{Path: "/task"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func newTask(c *gin.Context) {
	c.HTML(http.StatusOK, "addNew.html", gin.H{})
}
