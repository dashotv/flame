package models

type Connector struct {
	Download *DownloadStore
	Medium   *MediumStore
	Release  *ReleaseStore
}

func NewConnector() (*Connector, error) {
	download, err := NewDownloadStore("mongodb://localhost:27017", "seer_development", "downloads")
	if err != nil {
		return nil, err
	}
	medium, err := NewMediumStore("mongodb://localhost:27017", "seer_development", "media")
	if err != nil {
		return nil, err
	}
	release, err := NewReleaseStore("mongodb://localhost:27017", "torch_development", "torrents")
	if err != nil {
		return nil, err
	}

	c := &Connector{
		Download: download,
		Medium:   medium,
		Release:  release,
	}

	return c, nil
}
