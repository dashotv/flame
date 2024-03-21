package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /metube/
func (a *Application) MetubeIndex(c echo.Context) error {
	history, err := a.Metube.History()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": true, "message": "error loading Metube"})
	}
	return c.JSON(http.StatusOK, H{"error": false, "history": history})
}

// GET /metube/add
func (a *Application) MetubeAdd(c echo.Context, url string, name string) error {
	if name == "" {
		return c.JSON(http.StatusBadRequest, H{"error": true, "message": "name is required"})
	}
	if url == "" {
		return c.JSON(http.StatusBadRequest, H{"error": true, "message": "url is required"})
	}

	err := a.Metube.Add(name, url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": true, "message": "error adding to Metube: " + err.Error()})
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

// GET /metube/remove
func (a *Application) MetubeRemove(c echo.Context, name string, where string) error {
	if name == "" {
		return c.JSON(http.StatusBadRequest, H{"error": true, "message": "name is required"})
	}
	if where != "queue" && where != "done" {
		return c.JSON(http.StatusBadRequest, H{"error": true, "message": "where must be 'queue' or 'done'"})
	}

	err := a.Metube.Delete([]string{name}, where)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": true, "message": "error deleting from Metube: " + err.Error()})
	}

	return c.JSON(http.StatusOK, H{"error": false})
}
