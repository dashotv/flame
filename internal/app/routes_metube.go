package app

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /metube/
func (a *Application) MetubeIndex(c echo.Context) error {
	history, err := a.Metube.History()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "metube history: " + err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: history})
}

// GET /metube/add
func (a *Application) MetubeAdd(c echo.Context, url string, name string, autoStart bool) error {
	if name == "" {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: "name is required"})
	}
	if url == "" {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: "url is required"})
	}

	u, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		return err
	}

	app.Log.Named("metube").Debugf("add: %s %s", name, u)
	err = a.Metube.Add(name, string(u), autoStart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "metube add: " + err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}

// GET /metube/remove
func (a *Application) MetubeRemove(c echo.Context, name string, where string) error {
	if name == "" {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: "name is required"})
	}
	if where != "queue" && where != "done" {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: "where must be 'queue' or 'done'"})
	}

	err := a.Metube.Delete([]string{name}, where)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting from Metube: " + err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}
