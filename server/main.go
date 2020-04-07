package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dashotv/flame/server/nzbs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	nats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame/config"
	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/server/torrents"
	"github.com/dashotv/flame/utorrent"
	"github.com/dashotv/mercury"
)

type Server struct {
	cfg            *config.Config
	log            *logrus.Entry
	merc           *mercury.Mercury
	torrentChannel chan *utorrent.Response
	nzbChannel     chan *nzbget.GroupResponse
	torrent        *utorrent.Client
	nzb            *nzbget.Client
	cache          *redis.Client
}

func New(cfg *config.Config) (*Server, error) {
	var err error
	s := &Server{cfg: cfg}

	//if cfg.Mode == "dev" {
	//	logrus.SetLevel(logrus.DebugLevel)
	//}

	host, _ := os.Hostname()
	s.log = logrus.WithField("prefix", host)
	s.log.Level = logrus.DebugLevel

	s.merc, err = mercury.New("flame", nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "creating mercury")
	}

	s.torrentChannel = make(chan *utorrent.Response, 5)
	if err := s.merc.Sender("flame.torrents", s.torrentChannel); err != nil {
		return nil, errors.Wrap(err, "mercury sender")
	}

	s.cache = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // use default DB
	})

	s.log.Infof("configuring utorrent: %s", s.cfg.Utorrent.URL)
	s.torrent = utorrent.NewClient(s.cfg.Utorrent.URL)
	s.log.Infof("configuring nzbget: %s", s.cfg.Nzbget.URL)
	s.nzb = nzbget.NewClient(s.cfg.Nzbget.URL)

	return s, nil
}

func (s *Server) Start() error {
	s.log.Info("starting flame...")

	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("* * * * * *", s.SendTorrents); err != nil {
		return errors.Wrap(err, "adding cron function")
	}
	if _, err := c.AddFunc("* * * * * *", s.SendNzbs); err != nil {
		return errors.Wrap(err, "adding cron function")
	}

	go func() {
		s.log.Info("starting cron...")
		c.Start()
	}()

	if s.cfg.Mode == "release" {
		gin.SetMode(s.cfg.Mode)
	}
	router := gin.Default()
	router.GET("/", homeIndex)
	torrents.Routes(s.cache, s.torrent, router)
	nzbs.Routes(s.cache, s.nzb, router)

	s.log.Info("starting web...")
	if err := router.Run(fmt.Sprintf(":%d", s.cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (s *Server) SendTorrents() {
	resp, err := s.torrent.List()
	if err != nil {
		logrus.Errorf("couldn't get torrent list: %s", err)
		return
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		logrus.Errorf("couldn't marshal torrents: %s", err)
		return
	}

	s.cache.Set("flame-torrents", string(b), time.Minute)
	s.torrentChannel <- resp
}

func (s *Server) SendNzbs() {
	resp, err := s.nzb.List()
	if err != nil {
		logrus.Errorf("couldn't get nzb list: %s", err)
		return
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		logrus.Errorf("couldn't marshal nzbs: %s", err)
		return
	}

	s.cache.Set("flame-nzbs", string(b), time.Minute)
	s.nzbChannel <- resp
}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
