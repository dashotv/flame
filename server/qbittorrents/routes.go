// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package qbittorrents

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

	r := app.Router.Group("/qbittorrents")
	r.GET("/add", addHandler)
	r.GET("/", indexHandler)
	r.GET("/label", labelHandler)
	r.GET("/pause", pauseHandler)
	r.GET("/remove", removeHandler)
	r.GET("/resume", resumeHandler)
	r.GET("/want", wantHandler)
	r.GET("/wanted", wantedHandler)

}

func addHandler(c *gin.Context) {
	url := web.QueryString(c, "url")

	Add(c, url)
}

func indexHandler(c *gin.Context) {

	Index(c)
}

func labelHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")
	label := web.QueryString(c, "label")

	Label(c, infohash, label)
}

func pauseHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Pause(c, infohash)
}

func removeHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")
	del := web.QueryBool(c, "del")

	Remove(c, infohash, del)
}

func resumeHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Resume(c, infohash)
}

func wantHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")
	files := web.QueryString(c, "files")

	Want(c, infohash, files)
}

func wantedHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Wanted(c, infohash)
}
