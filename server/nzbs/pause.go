package nzbs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if id == -1 {
		PauseAll(c)
		return
	}

	err = client.Pause(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func PauseAll(c *gin.Context) {
	err := client.PauseAll()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Resume(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if id == -1 {
		ResumeAll(c)
		return
	}
	err = client.Resume(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ResumeAll(c *gin.Context) {
	err := client.ResumeAll()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
