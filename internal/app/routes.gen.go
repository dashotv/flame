// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/plugins/router"
	"github.com/labstack/echo/v4"
)

func init() {
	initializers = append(initializers, setupRoutes)
	healthchecks["routes"] = checkRoutes
	starters = append(starters, startRoutes)
}

func checkRoutes(app *Application) error {
	// TODO: check routes
	return nil
}

func startRoutes(ctx context.Context, app *Application) error {
	go func() {
		app.Routes()
		app.Log.Info("starting routes...")
		if err := app.Engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
			app.Log.Errorf("routes: %s", err)
		}
	}()
	return nil
}

func setupRoutes(app *Application) error {
	logger := app.Log.Named("routes").Desugar()
	e, err := router.New(logger)
	if err != nil {
		return fae.Wrap(err, "router plugin")
	}
	app.Engine = e
	// unauthenticated routes
	app.Default = app.Engine.Group("")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("")

	// TODO: fix auth
	if app.Config.Auth {
		clerkSecret := app.Config.ClerkSecretKey
		if clerkSecret == "" {
			app.Log.Fatal("CLERK_SECRET_KEY is not set")
		}
		clerkToken := app.Config.ClerkToken
		if clerkToken == "" {
			app.Log.Fatal("CLERK_TOKEN is not set")
		}

		app.Router.Use(router.ClerkAuth(clerkSecret, clerkToken))
	}

	return nil
}

type Setting struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type SettingsBatch struct {
	IDs   []string `json:"ids"`
	Name  string   `json:"name"`
	Value bool     `json:"value"`
}

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func (a *Application) Routes() {
	a.Default.GET("/", a.indexHandler)
	a.Default.GET("/health", a.healthHandler)

	metube := a.Router.Group("/metube")
	metube.GET("/", a.MetubeIndexHandler)
	metube.GET("/add", a.MetubeAddHandler)
	metube.GET("/remove", a.MetubeRemoveHandler)

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

func (a *Application) indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, router.H{
		"name": "flame",
		"routes": router.H{
			"metube":       "/metube",
			"nzbs":         "/nzbs",
			"qbittorrents": "/qbittorrents",
		},
	})
}

func (a *Application) healthHandler(c echo.Context) error {
	health, err := a.Health()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, router.H{"name": "flame", "health": health})
}

// Metube (/metube)
func (a *Application) MetubeIndexHandler(c echo.Context) error {
	return a.MetubeIndex(c)
}
func (a *Application) MetubeAddHandler(c echo.Context) error {
	url := router.QueryParamString(c, "url")
	name := router.QueryParamString(c, "name")
	auto_start := router.QueryParamBool(c, "auto_start")
	return a.MetubeAdd(c, url, name, auto_start)
}
func (a *Application) MetubeRemoveHandler(c echo.Context) error {
	name := router.QueryParamString(c, "name")
	where := router.QueryParamString(c, "where")
	return a.MetubeRemove(c, name, where)
}

// Nzbs (/nzbs)
func (a *Application) NzbsIndexHandler(c echo.Context) error {
	return a.NzbsIndex(c)
}
func (a *Application) NzbsAddHandler(c echo.Context) error {
	url := router.QueryParamString(c, "url")
	category := router.QueryParamString(c, "category")
	name := router.QueryParamString(c, "name")
	return a.NzbsAdd(c, url, category, name)
}
func (a *Application) NzbsRemoveHandler(c echo.Context) error {
	id := router.QueryParamInt(c, "id")
	return a.NzbsRemove(c, id)
}
func (a *Application) NzbsDestroyHandler(c echo.Context) error {
	id := router.QueryParamInt(c, "id")
	return a.NzbsDestroy(c, id)
}
func (a *Application) NzbsPauseHandler(c echo.Context) error {
	id := router.QueryParamInt(c, "id")
	return a.NzbsPause(c, id)
}
func (a *Application) NzbsResumeHandler(c echo.Context) error {
	id := router.QueryParamInt(c, "id")
	return a.NzbsResume(c, id)
}
func (a *Application) NzbsHistoryHandler(c echo.Context) error {
	hidden := router.QueryParamBool(c, "hidden")
	return a.NzbsHistory(c, hidden)
}

// Qbittorrents (/qbittorrents)
func (a *Application) QbittorrentsIndexHandler(c echo.Context) error {
	return a.QbittorrentsIndex(c)
}
func (a *Application) QbittorrentsAddHandler(c echo.Context) error {
	url := router.QueryParamString(c, "url")
	return a.QbittorrentsAdd(c, url)
}
func (a *Application) QbittorrentsRemoveHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	del := router.QueryParamBool(c, "del")
	return a.QbittorrentsRemove(c, infohash, del)
}
func (a *Application) QbittorrentsPauseHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	return a.QbittorrentsPause(c, infohash)
}
func (a *Application) QbittorrentsResumeHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	return a.QbittorrentsResume(c, infohash)
}
func (a *Application) QbittorrentsLabelHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	label := router.QueryParamString(c, "label")
	return a.QbittorrentsLabel(c, infohash, label)
}
func (a *Application) QbittorrentsWantHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	files := router.QueryParamString(c, "files")
	return a.QbittorrentsWant(c, infohash, files)
}
func (a *Application) QbittorrentsWantedHandler(c echo.Context) error {
	infohash := router.QueryParamString(c, "infohash")
	return a.QbittorrentsWanted(c, infohash)
}
