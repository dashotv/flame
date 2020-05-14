package nzbs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context, id int) {
	var err error

	if id == -1 {
		PauseAll(c)
		return
	}

	err = app.Nzbget.Pause(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func PauseAll(c *gin.Context) {
	err := app.Nzbget.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Resume(c *gin.Context, id int) {
	var err error

	if id == -1 {
		ResumeAll(c)
		return
	}
	err = app.Nzbget.Resume(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ResumeAll(c *gin.Context) {
	err := app.Nzbget.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
