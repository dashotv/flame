package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownloadsIndex(c *gin.Context) {
	results, err := App().DB.ActiveDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, d := range results {
		m := &Medium{}

		if err := App().DB.Medium.FindByID(d.MediumId, m); err != nil {
			App().Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		App().Log.Infof("found %s: %s", m.ID, m.Title)
	}

	c.JSON(http.StatusOK, results)
}
