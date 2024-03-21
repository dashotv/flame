package app

import (
	"github.com/dashotv/flame/nzbget"
)

func init() {
	initializers = append(initializers, setupNzbget)
}

func setupNzbget(app *Application) error {
	app.Log.Infof("connecting nzbget: %s", app.Config.NzbgetURL)
	app.Nzb = nzbget.NewClient(app.Config.NzbgetURL)
	return nil
}
