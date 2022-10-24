// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() {
	s.Router.GET("/", homeHandler)

	downloads := s.Router.Group("/downloads")
	downloads.GET("/", downloadsIndexHandler)

	media := s.Router.Group("/media")
	media.GET("/", mediaIndexHandler)

	nzbs := s.Router.Group("/nzbs")
	nzbs.GET("/add", nzbsAddHandler)
	nzbs.GET("/destroy", nzbsDestroyHandler)
	nzbs.GET("/history", nzbsHistoryHandler)
	nzbs.GET("/", nzbsIndexHandler)
	nzbs.GET("/pause", nzbsPauseHandler)
	nzbs.GET("/remove", nzbsRemoveHandler)
	nzbs.GET("/resume", nzbsResumeHandler)

	qbittorrents := s.Router.Group("/qbittorrents")
	qbittorrents.GET("/add", qbittorrentsAddHandler)
	qbittorrents.GET("/", qbittorrentsIndexHandler)
	qbittorrents.GET("/label", qbittorrentsLabelHandler)
	qbittorrents.GET("/pause", qbittorrentsPauseHandler)
	qbittorrents.GET("/remove", qbittorrentsRemoveHandler)
	qbittorrents.GET("/resume", qbittorrentsResumeHandler)
	qbittorrents.GET("/want", qbittorrentsWantHandler)
	qbittorrents.GET("/wanted", qbittorrentsWantedHandler)

}

func homeHandler(c *gin.Context) {
	Index(c)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "home")
}

// /downloads
func downloadsIndexHandler(c *gin.Context) {

	DownloadsIndex(c)
}

// /media
func mediaIndexHandler(c *gin.Context) {

	MediaIndex(c)
}

// /nzbs
func nzbsAddHandler(c *gin.Context) {
	url := c.Param("url")
	category := c.Param("category")
	name := c.Param("name")

	NzbsAdd(c, url, category, name)
}

func nzbsDestroyHandler(c *gin.Context) {
	id := c.Param("id")

	NzbsDestroy(c, id)
}

func nzbsHistoryHandler(c *gin.Context) {
	hidden := c.Param("hidden")

	NzbsHistory(c, hidden)
}

func nzbsIndexHandler(c *gin.Context) {

	NzbsIndex(c)
}

func nzbsPauseHandler(c *gin.Context) {
	id := c.Param("id")

	NzbsPause(c, id)
}

func nzbsRemoveHandler(c *gin.Context) {
	id := c.Param("id")

	NzbsRemove(c, id)
}

func nzbsResumeHandler(c *gin.Context) {
	id := c.Param("id")

	NzbsResume(c, id)
}

// /qbittorrents
func qbittorrentsAddHandler(c *gin.Context) {
	url := c.Param("url")

	QbittorrentsAdd(c, url)
}

func qbittorrentsIndexHandler(c *gin.Context) {

	QbittorrentsIndex(c)
}

func qbittorrentsLabelHandler(c *gin.Context) {
	infohash := c.Param("infohash")
	label := c.Param("label")

	QbittorrentsLabel(c, infohash, label)
}

func qbittorrentsPauseHandler(c *gin.Context) {
	infohash := c.Param("infohash")

	QbittorrentsPause(c, infohash)
}

func qbittorrentsRemoveHandler(c *gin.Context) {
	infohash := c.Param("infohash")
	del := c.Param("del")

	QbittorrentsRemove(c, infohash, del)
}

func qbittorrentsResumeHandler(c *gin.Context) {
	infohash := c.Param("infohash")

	QbittorrentsResume(c, infohash)
}

func qbittorrentsWantHandler(c *gin.Context) {
	infohash := c.Param("infohash")
	files := c.Param("files")

	QbittorrentsWant(c, infohash, files)
}

func qbittorrentsWantedHandler(c *gin.Context) {
	infohash := c.Param("infohash")

	QbittorrentsWanted(c, infohash)
}
