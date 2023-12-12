package app

import (
	"context"
	"fmt"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/dashotv/mercury"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

var server *Server
var ctx = context.Background()

func setupServer() (err error) {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	server = &Server{
		Log: log.Named("server"),
	}

	server.Engine = gin.New()
	server.Engine.Use(ginzap.Ginzap(log.Desugar(), time.RFC3339, true), ginzap.RecoveryWithZap(log.Desugar(), true))
	server.Default = server.Engine.Group("/")
	server.Router = server.Engine.Group("/")

	server.Routes()

	// 	if cfg.Auth {
	// 		clerkSecret := cfg.ClerkSecretKey
	// 		if clerkSecret == "" {
	// 			log.Fatal("CLERK_SECRET_KEY is not set")
	// 		}
	//
	// 		clerkClient, err := clerk.NewClient(clerkSecret)
	// 		if err != nil {
	// 			log.Fatalf("clerk: %s", err)
	// 		}
	//
	// 		server.Router.Use(requireSession(clerkClient))
	// 	}

	server.merc, err = mercury.New("flame", cfg.NatsURL)
	if err != nil {
		return errors.Wrap(err, "creating mercury")
	}

	server.combined = make(chan *Combined, 5)
	if err := server.merc.Sender("flame.combined", server.combined); err != nil {
		return errors.Wrap(err, "mercury sender")
	}

	server.qbtChannel = make(chan *qbt.Response, 5)
	if err := server.merc.Sender("flame.qbittorrents", server.qbtChannel); err != nil {
		return errors.Wrap(err, "mercury sender")
	}

	server.nzbChannel = make(chan *nzbget.GroupResponse, 5)
	if err := server.merc.Sender("flame.nzbs", server.nzbChannel); err != nil {
		return errors.Wrap(err, "mercury sender")
	}

	server.metricsChannel = make(chan *Metrics, 5)
	if err := server.merc.Sender("flame.metrics", server.metricsChannel); err != nil {
		return errors.Wrap(err, "mercury sender")
	}

	return nil
}

type Server struct {
	Engine  *gin.Engine
	Router  *gin.RouterGroup
	Default *gin.RouterGroup
	Log     *zap.SugaredLogger

	merc           *mercury.Mercury
	combined       chan *Combined
	qbtChannel     chan *qbt.Response
	nzbChannel     chan *nzbget.GroupResponse
	metricsChannel chan *Metrics
}

func (s *Server) Start() error {
	s.Log.Info("starting flame...")

	if cfg.Cron {
		c := cron.New(cron.WithSeconds())

		if _, err := c.AddFunc("* * * * * *", s.Updates); err != nil {
			return errors.Wrap(err, "adding updates cron function")
		}

		go func() {
			s.Log.Info("starting cron...")
			c.Start()
		}()
	}

	s.Log.Info("starting web...")
	if err := s.Engine.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

// pauseAll pauses all torrents and sets a flag in the cache
func (s *Server) pauseAll() error {
	s.Log.Info("pausing all...")
	err := qb.PauseAll()
	if err != nil {
		return errors.Wrap(err, "pausing all")
	}

	status := cache.Set(ctx, "flame-disk-paused", true, time.Hour*24*7)
	if status.Err() != nil {
		s.Log.Errorf("sendqbts: set cache failed: %s", status.Err())
	}

	return nil
}

func (s *Server) resumeAll() error {
	err := qb.ResumeAll()
	if err != nil {
		s.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}

	status := cache.Del(ctx, "flame-disk-paused")
	if status.Err() != nil {
		s.Log.Errorf("sendqbts: set cache failed: %s", status.Err())
	}

	return nil
}

// diskPaused checks the cache for the flag
func (s *Server) diskPaused() bool {
	paused, err := cache.Get(ctx, "flame-disk-paused").Result()
	if err != nil {
		s.Log.Errorf("paused: %s", err)
		return false
	}
	return paused == "true"
}

// allPaused checks if all torrents are paused and the paused flag is true
func (s *Server) allPaused() bool {
	if !s.diskPaused() {
		return false
	}

	ok, err := qb.AllPaused()
	if err != nil {
		s.Log.Errorf("checkdisk: failed to check if all qbts are paused: %s", err)
	}
	return ok
}

func (s *Server) checkDisk(resp *nzbget.GroupResponse, qbt *qbt.Response) {
	if resp.Status.FreeDiskSpaceMB < 25000 && !s.allPaused() {
		err := s.pauseAll()
		if err != nil {
			s.Log.Errorf("checkdisk: failed to pause all qbts: %s", err)
		}
		return
	}

	if !s.diskPaused() {
		return
	}
	if !qbt.AllPaused() {
		return
	}

	err := s.resumeAll()
	if err != nil {
		s.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}

	return
}
