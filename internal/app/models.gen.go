// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"go.uber.org/zap"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/grimoire"
	"github.com/kamva/mgm/v3"
)

func init() {
	initializers = append(initializers, setupDb)
	healthchecks["db"] = checkDb
}

func setupDb(app *Application) error {
	db, err := NewConnector(app)
	if err != nil {
		return err
	}

	app.DB = db
	return nil
}

func checkDb(app *Application) (err error) {
	// TODO: Check DB connection
	return nil
}

type Connector struct {
	Log *zap.SugaredLogger
}

func connection[T mgm.Model](name string) (*grimoire.Store[T], error) {
	s, err := app.Config.ConnectionFor(name)
	if err != nil {
		return nil, err
	}
	c, err := grimoire.New[T](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewConnector(app *Application) (*Connector, error) {

	c := &Connector{
		Log: app.Log.Named("db"),
	}

	return c, nil
}

type Combined struct { // struct
	Torrents  []*qbt.Torrent `bson:"torrents" json:"torrents"`
	Nzbs      []nzbget.Group `bson:"nzbs" json:"nzbs"`
	NzbStatus nzbget.Status  `bson:"nzb_status" json:"nzb_status"`
	Metrics   *Metrics       `bson:"metrics" json:"metrics"`
}

type Metrics struct { // struct
	Diskspace string `bson:"diskspace" json:"diskspace"`
	Torrents  struct {
		DownloadRate string `json:"download_rate"`
		UploadRate   string `json:"upload_rate"`
	} `bson:"torrents" json:"torrents"`
	Nzbs struct {
		DownloadRate string `json:"download_rate"`
	} `bson:"nzbs" json:"nzbs"`
}
