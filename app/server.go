package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/mercury"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

var ctx = context.Background()

type Server struct {
	Router *gin.Engine
	Log    *logrus.Entry
	App    *Application
	Config *Config

	merc           *mercury.Mercury
	qbtChannel     chan *qbt.Response
	nzbChannel     chan *nzbget.GroupResponse
	metricsChannel chan *Metrics
}

func New() (*Server, error) {
	var err error
	s := &Server{
		Log:    App().Log,
		Router: App().Router,
		Config: ConfigInstance(),
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

	s.metricsChannel = make(chan *Metrics, 5)
	if err := s.merc.Sender("flame.metrics", s.metricsChannel); err != nil {
		return nil, errors.Wrap(err, "mercury sender")
	}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting flame...")

	if s.Config.Cron {
		c := cron.New(cron.WithSeconds())

		if _, err := c.AddFunc("* * * * * *", s.Updates); err != nil {
			return errors.Wrap(err, "adding updates cron function")
		}

		go func() {
			s.Log.Info("starting cron...")
			c.Start()
		}()
	}

	s.Routes()

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (s *Server) checkDisk(resp *nzbget.GroupResponse) {
	//s.Log.Infof("checkdisk: checking free disk space: %d MB", resp.Status.FreeDiskSpaceMB)
	if resp.Status.FreeDiskSpaceMB < 25000 {
		s.Log.Warnf("checkdisk: free disk space low")
		err := App().Qbittorrent.PauseAll()
		if err != nil {
			s.Log.Errorf("checkdisk: failed to pause all qbts: %s", err)
		}
	} else {
		ok, err := App().Qbittorrent.AllPaused()
		if err != nil {
			s.Log.Errorf("checkdisk: failed to check if all qbts are paused: %s", err)
		}
		if ok {
			s.Log.Infof("checkdisk: free disk space restored")
			err := App().Qbittorrent.ResumeAll()
			if err != nil {
				s.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
			}
		}
	}
}
