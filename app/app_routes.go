// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers = append(initializers, setupRoutes)
	healthchecks["routes"] = checkRoutes
}

func checkRoutes(app *Application) error {
	// TODO: check routes
	return nil
}

func setupRoutes(app *Application) error {
	if app.Config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := app.Log.Named("routes").Desugar()

	app.Engine = gin.New()
	app.Engine.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)
	// unauthenticated routes
	app.Default = app.Engine.Group("/")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("/")

	// if app.Config.Auth {
	// 	clerkSecret := app.Config.ClerkSecretKey
	// 	if clerkSecret == "" {
	// 		app.Log.Fatal("CLERK_SECRET_KEY is not set")
	// 	}
	//
	// 	clerkClient, err := clerk.NewClient(clerkSecret)
	// 	if err != nil {
	// 		app.Log.Fatalf("clerk: %s", err)
	// 	}
	//
	// 	app.Router.Use(requireSession(clerkClient))
	// }

	return nil
}

// Enable Auth and uncomment to use Clerk to manage auth
// also add this import: "github.com/clerkinc/clerk-sdk-go/clerk"
//
// requireSession wraps the clerk.RequireSession middleware
// func requireSession(client clerk.Client) gin.HandlerFunc {
// 	requireActiveSession := clerk.RequireSessionV2(client)
// 	return func(gctx *gin.Context) {
// 		var skip = true
// 		var handler http.HandlerFunc = func(http.ResponseWriter, *http.Request) {
// 			skip = false
// 		}
// 		requireActiveSession(handler).ServeHTTP(gctx.Writer, gctx.Request)
// 		switch {
// 		case skip:
// 			gctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session required"})
// 		default:
// 			gctx.Next()
// 		}
// 	}
// }

func (a *Application) Routes() {
	a.Default.GET("/", a.indexHandler)
	a.Default.GET("/health", a.healthHandler)

	nzbs := a.Router.Group("/nzbs")
	nzbs.GET("/", a.NzbsIndexHandler)
	nzbs.GET("/add", a.NzbsAddHandler)
	nzbs.GET("/remove", a.NzbsRemoveHandler)
	nzbs.GET("/destroy", a.NzbsDestroyHandler)
	nzbs.GET("/pause", a.NzbsPauseHandler)
	nzbs.GET("/resume", a.NzbsResumeHandler)
	nzbs.GET("/history", a.NzbsHistoryHandler)

	qbittorrents := a.Router.Group("/qbittorrents")
	qbittorrents.GET("/", a.QbittorrentsIndexHandler)
	qbittorrents.GET("/add", a.QbittorrentsAddHandler)
	qbittorrents.GET("/remove", a.QbittorrentsRemoveHandler)
	qbittorrents.GET("/pause", a.QbittorrentsPauseHandler)
	qbittorrents.GET("/resume", a.QbittorrentsResumeHandler)
	qbittorrents.GET("/label", a.QbittorrentsLabelHandler)
	qbittorrents.GET("/want", a.QbittorrentsWantHandler)
	qbittorrents.GET("/wanted", a.QbittorrentsWantedHandler)

}

func (a *Application) indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "flame",
		"routes": gin.H{
			"nzbs":         "/nzbs",
			"qbittorrents": "/qbittorrents",
		},
	})
}

func (a *Application) healthHandler(c *gin.Context) {
	health, err := a.Health()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": "flame", "health": health})
}

// Nzbs (/nzbs)
func (a *Application) NzbsIndexHandler(c *gin.Context) {
	a.NzbsIndex(c)
}
func (a *Application) NzbsAddHandler(c *gin.Context) {
	url := QueryString(c, "url")
	category := QueryString(c, "category")
	name := QueryString(c, "name")
	a.NzbsAdd(c, url, category, name)
}
func (a *Application) NzbsRemoveHandler(c *gin.Context) {
	id := QueryInt(c, "id")
	a.NzbsRemove(c, id)
}
func (a *Application) NzbsDestroyHandler(c *gin.Context) {
	id := QueryInt(c, "id")
	a.NzbsDestroy(c, id)
}
func (a *Application) NzbsPauseHandler(c *gin.Context) {
	id := QueryInt(c, "id")
	a.NzbsPause(c, id)
}
func (a *Application) NzbsResumeHandler(c *gin.Context) {
	id := QueryInt(c, "id")
	a.NzbsResume(c, id)
}
func (a *Application) NzbsHistoryHandler(c *gin.Context) {
	hidden := QueryBool(c, "hidden")
	a.NzbsHistory(c, hidden)
}

// Qbittorrents (/qbittorrents)
func (a *Application) QbittorrentsIndexHandler(c *gin.Context) {
	a.QbittorrentsIndex(c)
}
func (a *Application) QbittorrentsAddHandler(c *gin.Context) {
	url := QueryString(c, "url")
	a.QbittorrentsAdd(c, url)
}
func (a *Application) QbittorrentsRemoveHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	del := QueryBool(c, "del")
	a.QbittorrentsRemove(c, infohash, del)
}
func (a *Application) QbittorrentsPauseHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	a.QbittorrentsPause(c, infohash)
}
func (a *Application) QbittorrentsResumeHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	a.QbittorrentsResume(c, infohash)
}
func (a *Application) QbittorrentsLabelHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	label := QueryString(c, "label")
	a.QbittorrentsLabel(c, infohash, label)
}
func (a *Application) QbittorrentsWantHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	files := QueryString(c, "files")
	a.QbittorrentsWant(c, infohash, files)
}
func (a *Application) QbittorrentsWantedHandler(c *gin.Context) {
	infohash := QueryString(c, "infohash")
	a.QbittorrentsWanted(c, infohash)
}