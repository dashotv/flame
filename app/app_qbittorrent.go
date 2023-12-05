package app

import (
	"github.com/pkg/errors"

	"github.com/dashotv/flame/qbt"
)

var qb *qbt.Api

func setupQbittorrent() error {
	log.Infof("connecting qbittorrent: %s", cfg.QbittorrentURL)
	qb = qbt.NewApi(cfg.QbittorrentURL)
	ok, err := qb.Login(cfg.QbittorrentUsername, cfg.QbittorrentPassword)
	if err != nil {
		return errors.Errorf("qbittorrent: could not login: %s", err)
	}
	if !ok {
		return errors.Errorf("qbittorrent: login false")
	}
	return nil
}
