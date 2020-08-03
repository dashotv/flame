package torrents

//func Start(c *gin.Context, infohash string) {
//	err := app.Qbittorrent.Start(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}
//
//func Stop(c *gin.Context, infohash string) {
//	err := app.Utorrent.Stop(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}
