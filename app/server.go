package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/mercury"
)

var ctx = context.Background()

type Server struct {
	Router *gin.Engine
	Log    *logrus.Entry
	App    *Application
	Config *Config

	merc       *mercury.Mercury
	qbtChannel chan *qbt.Response
	nzbChannel chan *nzbget.GroupResponse
}

func New() (*Server, error) {
	var err error
	cfg := ConfigInstance()
	app := App()
	fmt.Printf("%v\n", app)
	log := app.Log.WithField("prefix", "server")
	s := &Server{
		Log:    log,
		Router: app.Router,
		App:    app,
		Config: cfg,
	}

	s.merc, err = mercury.New("flame", s.Config.Nats.URL)
	if err != nil {
		return nil, errors.Wrap(err, "creating mercury")
	}

	s.qbtChannel = make(chan *qbt.Response, 5)
	if err := s.merc.Sender("flame.qbittorrents", s.qbtChannel); err != nil {
		return nil, errors.Wrap(err, "mercury sender")
	}

	s.nzbChannel = make(chan *nzbget.GroupResponse, 5)
	if err := s.merc.Sender("flame.nzbs", s.nzbChannel); err != nil {
		return nil, errors.Wrap(err, "mercury sender")
	}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting flame...")

	c := cron.New(cron.WithSeconds())
	if s.Config.Cron {
		// every second SendQbittorrents
		if _, err := c.AddFunc("* * * * * *", s.SendQbittorrents); err != nil {
			return errors.Wrap(err, "adding cron function")
		}
		// every second SendQbittorrents
		if _, err := c.AddFunc("* * * * * *", s.SendNzbs); err != nil {
			return errors.Wrap(err, "adding cron function")
		}
	}

	go func() {
		s.Log.Info("starting cron...")
		c.Start()
	}()

	s.Routes()

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (s *Server) SendQbittorrents() {
	resp, err := s.App.Qbittorrent.List()
	if err != nil {
		logrus.Errorf("couldn't get torrent list: %s", err)
		return
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		logrus.Errorf("couldn't marshal torrents: %s", err)
		return
	}

	s.App.Cache.Set(ctx, "flame-qbittorrents", string(b), time.Minute)
	s.qbtChannel <- resp
}

func (s *Server) SendNzbs() {
	resp, err := s.App.Nzbget.List()
	if err != nil {
		logrus.Errorf("couldn't get nzb list: %s", err)
		return
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		logrus.Errorf("couldn't marshal nzbs: %s", err)
		return
	}

	s.App.Cache.Set(ctx, "flame-nzbs", string(b), time.Minute)
	s.nzbChannel <- resp
}
