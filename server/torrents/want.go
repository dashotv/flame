package torrents

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Want(c *gin.Context) {
	infohash := c.Query("infohash")
	files := c.Query("files")
	if files == "none" {
		WantNone(c, infohash)
		return
	}
	ids, err := filesToIds(files)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = client.Want(infohash, ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func filesToIds(files string) ([]int, error) {
	list := strings.Split(files, ",")
	ids := make([]int, len(list))
	for i, v := range list {
		num, err := strconv.Atoi(v)
		if err != nil {
			return ids, err
		}
		ids[i] = num
	}
	return ids, nil
}

func WantNone(c *gin.Context, infohash string) {
	err := client.WantNone(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Wanted(c *gin.Context) {
	infohash := c.Query("infohash")
	want, err := client.Wanted(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "wanted": want})
}
