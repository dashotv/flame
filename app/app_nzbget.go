package app

import "github.com/dashotv/flame/nzbget"

var nzb *nzbget.Client

func setupNzbget() error {
	log.Infof("connecting nzbget: %s", cfg.NzbgetURL)
	nzb = nzbget.NewClient(cfg.NzbgetURL)
	return nil
}
