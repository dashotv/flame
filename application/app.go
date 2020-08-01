package application

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/dashotv/flame/config"
	"github.com/dashotv/flame/models"
	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/flame/utorrent"
)

var once sync.Once
var instance *App

type App struct {
	Config *config.Config
	Router *gin.Engine
	Cache  *redis.Client
	Log    *logrus.Entry
	DB     *models.Connector

	Nzbget      *nzbget.Client
	Utorrent    *utorrent.Client
	Qbittorrent *qbt.Client
}

func logger() *logrus.Entry {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&prefixed.TextFormatter{DisableTimestamp: false, FullTimestamp: true})
	host, _ := os.Hostname()
	return logrus.WithField("prefix", host)
}

func initialize() *App {
	cfg := config.Instance()
	log := logger()

	db, err := models.NewConnector()
	if err != nil {
		log.Errorf("database connection failed: %s", err)
	}

	if cfg.Mode == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cfg.Mode == "release" {
		gin.SetMode(cfg.Mode)
	}

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// TODO: add this to config
	cache := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // use default DB
	})

	log.Infof("connecting nzbget: %s", cfg.Nzbget.URL)
	nzb := nzbget.NewClient(cfg.Nzbget.URL)

	log.Infof("connecting utorrent: %s", cfg.Utorrent.URL)
	ut := utorrent.NewClient(cfg.Utorrent.URL)

	log.Infof("connecting qbittorrent: %s", cfg.Qbittorrent.URL)
	qb := qbt.NewClient(cfg.Qbittorrent.URL)
	ok, err := qb.Login(cfg.Qbittorrent.Username, cfg.Qbittorrent.Password)
	if err != nil {
		log.Errorf("qbittorrent: could not login: %s", err)
	}
	if !ok {
		log.Errorf("qbittorrent: login false")
	}

	return &App{
		Config:      cfg,
		DB:          db,
		Nzbget:      nzb,
		Utorrent:    ut,
		Qbittorrent: qb,
		Router:      router,
		Cache:       cache,
		Log:         log,
	}
}

func Instance() *App {
	once.Do(func() {
		instance = initialize()
	})
	return instance
}
