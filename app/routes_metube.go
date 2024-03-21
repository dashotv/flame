package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /metube/
func (a *Application) MetubeIndex(c echo.Context) error {
	// list, err := a.DB.MetubeList()
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, H{"error": true, "message": "error loading Metube"})
	// }
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}

// GET /metube/add
func (a *Application) MetubeAdd(c echo.Context, url string, name string) error {
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}

// GET /metube/remove
func (a *Application) MetubeRemove(c echo.Context, name string) error {
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}
