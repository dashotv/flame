package app

import (
	"github.com/pkg/errors"

	"github.com/dashotv/flame/qbt"
)

func init() {
	initializers = append(initializers, setupQbittorrent)
}

func setupQbittorrent(app *Application) error {
	app.Log.Infof("connecting qbittorrent: %s", app.Config.QbittorrentURL)
	app.Qbt = qbt.NewApi(app.Config.QbittorrentURL)
	ok, err := app.Qbt.Login(app.Config.QbittorrentUsername, app.Config.QbittorrentPassword)
	if err != nil {
		return errors.Errorf("qbittorrent: could not login: %s", err)
	}
	if !ok {
		return errors.Errorf("qbittorrent: login false")
	}
	return nil
}
