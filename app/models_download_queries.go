package app

func (c *Connector) ActiveDownloads() ([]*Download, error) {
	q := c.Download.Query()
	return q.In("status", []string{"searching", "loading", "managing", "downloading", "reviewing"}).Run()
}
