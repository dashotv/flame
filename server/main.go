package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	nats "github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame/server/torrents"
	"github.com/dashotv/flame/utorrent"
	"github.com/dashotv/mercury"
)

type Server struct {
	URL  string
	Port int
	Mode string

	log     *logrus.Entry
	merc    *mercury.Mercury
	channel chan *utorrent.Response
	client  *utorrent.Client
	cache   *redis.Client
}

func New(URL, mode string, port int) (*Server, error) {
	var err error
	s := &Server{URL: URL, Mode: mode, Port: port}

	host, _ := os.Hostname()
	s.log = logrus.WithField("prefix", host)

	s.merc, err = mercury.New("flame", nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "creating mercury")
	}

	s.channel = make(chan *utorrent.Response, 5)
	if err := s.merc.Sender("flame.torrents", s.channel); err != nil {
		return nil, errors.Wrap(err, "mercury sender")
	}

	s.cache = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // use default DB
	})

	s.client = utorrent.NewClient(URL)

	return s, nil
}

func (s *Server) Start() error {
	s.log.Info("starting flame...")

	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("* * * * * *", s.Sender); err != nil {
		return errors.Wrap(err, "adding cron function")
	}

	go func() {
		s.log.Info("starting cron...")
		c.Start()
	}()

	if s.Mode == "release" {
		gin.SetMode(s.Mode)
	}
	router := gin.Default()
	router.GET("/", homeIndex)
	torrents.Routes(s.cache, router)

	s.log.Info("starting web...")
	if err := router.Run(fmt.Sprintf(":%d", s.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (s *Server) Sender() {
	resp, err := s.client.List()
	if err != nil {
		logrus.Errorf("flame list error: %s", err)
		return
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		logrus.Errorf("json marshall error: %s", err)
		return
	}

	s.cache.Set("flame", string(b), time.Minute)
	s.channel <- resp
}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
