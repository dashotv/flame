package app

import "time"

func (c *Connector) Upcoming() ([]*Medium, error) {
	q := c.Medium.Query()
	return q.Where("missing", nil).LessThanEqual("release_date", time.Now()).Run()
}
