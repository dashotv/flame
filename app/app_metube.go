package app

import "github.com/dashotv/flame/metube"

func init() {
	initializers = append(initializers, setupMetube)
}

func setupMetube(app *Application) error {
	app.Log.Infof("connecting metube: %s", app.Config.MetubeURL)
	app.Metube = metube.New(app.Config.MetubeURL, false)
	return nil
}
