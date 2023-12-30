package app

func (c *Connector) MinionGet(id string) (*Minion, error) {
	m, err := c.Minion.Get(id, &Minion{})
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) MinionList() ([]*Minion, error) {
	list, err := c.Minion.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
