package app

import (
	"fmt"

	"github.com/dashotv/flame/internal/metube"
)

func init() {
	initializers = append(initializers, setupMetube)
}

func setupMetube(app *Application) error {
	if app.Config.MetubeURL == "" {
		return fmt.Errorf("Metube URL is required: %s", app.Config.MetubeURL)
	}
	app.Log.Debugf("connecting metube: %s", app.Config.MetubeURL)
	app.Metube = metube.New(app.Config.MetubeURL, false)
	return nil
}
