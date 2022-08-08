package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	results, err := app.DB.Medium.Upcoming()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//for _, d := range results {
	//	m, err := app.DB.Medium.FindByID(d.MediumId)
	//	if err != nil {
	//		app.Log.Errorf("could not find medium: %s", d.MediumId)
	//		continue
	//	}
	//	app.Log.Infof("found %s: %s", m.ID, m.Title)
	//}

	c.JSON(http.StatusOK, results)
}
