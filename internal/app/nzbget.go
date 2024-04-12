package app

import (
	"github.com/dashotv/flame/internal/nzbget"
)

func init() {
	initializers = append(initializers, setupNzbget)
}

func setupNzbget(app *Application) error {
	app.Log.Debugf("connecting nzbget: %s", app.Config.NzbgetURL)
	app.Nzb = nzbget.NewClient(app.Config.NzbgetURL)
	return nil
}
