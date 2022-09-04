package downloads

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/models"
)

func Index(c *gin.Context) {
	results, err := app.DB.ActiveDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, d := range results {
		m := &models.Medium{}

		if err := app.DB.Medium.FindByID(d.MediumId, m); err != nil {
			app.Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		app.Log.Infof("found %s: %s", m.ID, m.Title)
	}

	c.JSON(http.StatusOK, results)
}
