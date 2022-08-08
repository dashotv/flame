// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/server/downloads"
	"github.com/dashotv/flame/server/nzbs"
	"github.com/dashotv/flame/server/qbittorrents"
)

func (s *Server) Routes() {
	s.Router.GET("/", homeHandler)

	downloads.Routes()
	nzbs.Routes()
	qbittorrents.Routes()

}

func homeHandler(c *gin.Context) {
	Home(c)
}

func Home(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
