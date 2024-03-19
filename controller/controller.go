package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/udvarid/task-manager-golang/authenticator"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/service"
)

var (
	activeConfiguration = &model.Configuration{}
)

type GetSession struct {
	Id string `json:"id"`
	Qr bool   `json:"qr"`
}

type GetIdAndSession struct {
	Id      string `json:"id"`
	Session string `json:"session"`
}

func Init(config *model.Configuration) {
	activeConfiguration = config
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", startPage)
	router.GET("/qr/:id/:session", getQr)
	router.GET("/qrvalid/:id/:session", qrvalid)
	router.POST("/validate/", validate)
	router.POST("/validateqr/", validateqr)
	router.GET("/checkin/:id/:session", checkInTask)
	router.GET("/task/", taskPage)
	router.POST("/delete/:delete_id", deleteTask)
	router.POST("/prolong/:id/:prolong_days", prolongTask)
	router.GET("/newTask/", newTask)
	router.POST("/addTask/", addTask)
	router.Run()
}

// TODO
// 1, QR code possibility at login? https://pkg.go.dev/github.com/skip2/go-qrcode#section-readme
//   - startpageen legyen egy qr kód kérő checkbox
//   - validate fv a checkboxtól függően vagy normálisan viselkedik, vagy meghívja a QR-es oldalt
//   - QR-es oldal a kapott ID alapján sessiont generál, ebből és az ID-ból egy QR kódot, ezt kirajzolja és meghívja a validate fv-t
//   - az authenticator.Validate ezesetben ne küldjön üzenetet
// 2, embeded-el a html templateket, akár ez alapján https://stackoverflow.com/questions/74975426/load-html-code-into-gin-framework-template

func startPage(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", "localhost", false, true)
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.HTML(http.StatusOK, "start.html", gin.H{
		"title": "Home Page",
	})
}

func qrvalid(c *gin.Context) {
	id := c.Param("id")
	session := c.Param("session")
	c.HTML(http.StatusOK, "qr.html", gin.H{
		"title":   "Home Page",
		"id":      id,
		"session": session,
	})
}

func getQr(c *gin.Context) {
	id := c.Param("id")
	session := c.Param("session")
	qrPath := activeConfiguration.RemoteAddress + "checkin/" + id + "/" + session
	png, _ := qrcode.Encode(qrPath, qrcode.Highest, 256)
	c.Data(http.StatusOK, "string", png)
}

func validate(c *gin.Context) {
	var getSession GetSession
	c.BindJSON(&getSession)
	newSession, err := authenticator.GiveSession(getSession.Id)
	if err != nil {
		redirectTo(c, "/")
		return
	}
	if getSession.Qr {
		redirectTo(c, "/qrvalid/"+getSession.Id+"/"+newSession)
		return
	}

	isValidatedInTime := authenticator.Validate(activeConfiguration, getSession.Id, newSession, true)

	if isValidatedInTime {
		c.SetCookie("id", getSession.Id, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		c.SetCookie("session", newSession, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		redirectTo(c, "/task")
	} else {
		redirectTo(c, "/")
	}

}

func validateqr(c *gin.Context) {
	var getIdAndSession GetIdAndSession
	c.BindJSON(&getIdAndSession)

	isValidatedInTime := authenticator.Validate(activeConfiguration, getIdAndSession.Id, getIdAndSession.Session, false)

	if isValidatedInTime {
		c.SetCookie("id", getIdAndSession.Id, 3600, "/", activeConfiguration.RemoteAddress, false, true)
		c.SetCookie("session", getIdAndSession.Session, 3600, "/", activeConfiguration.RemoteAddress, false, true)
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
