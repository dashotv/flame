// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package nzbs

import (
	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/application"
	"github.com/dashotv/flame/config"
	"github.com/dashotv/golem/web"
)

var cfg *config.Config
var app *application.App

func Routes() {
	cfg = config.Instance()
	app = application.Instance()

	r := app.Router.Group("/nzbs")
	r.GET("/add", addHandler)
	r.GET("/destroy", destroyHandler)
	r.GET("/history", historyHandler)
	r.GET("/", indexHandler)
	r.GET("/pause", pauseHandler)
	r.GET("/remove", removeHandler)
	r.GET("/resume", resumeHandler)

}

func addHandler(c *gin.Context) {
	url := web.QueryString(c, "url")
	category := web.QueryString(c, "category")
	name := web.QueryString(c, "name")

	Add(c, url, category, name)
}

func destroyHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	Destroy(c, id)
}

func historyHandler(c *gin.Context) {
	hidden := web.QueryBool(c, "hidden")

	History(c, hidden)
}

func indexHandler(c *gin.Context) {

	Index(c)
}

func pauseHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	Pause(c, id)
}

func removeHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	Remove(c, id)
}

func resumeHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	Resume(c, id)
}
