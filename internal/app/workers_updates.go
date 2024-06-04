package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sony/gobreaker/v2"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

var qbtCircuitBreaker *gobreaker.CircuitBreaker[*qbt.Response]

func init() {
	initializers = append(initializers, setupUpdates)
}

func setupUpdates(a *Application) error {
	qbtCircuitBreaker = gobreaker.NewCircuitBreaker[*qbt.Response](gobreaker.Settings{
		Name: "qbt",
		// ReadyToTrip: func(counts gobreaker.Counts) bool {
		// 	return counts.ConsecutiveFailures > 3
		// },
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			a.Workers.Log.Infof("qbt: state change: %s -> %s", from, to)
			if to == gobreaker.StateHalfOpen {
				if err := setupQbittorrent(a); err != nil {
					a.Workers.Log.Errorf("qbt: re-init failed: %s", err)
				}
			}
		},
		Timeout: 5 * time.Second,
	})
	return nil
}

func Updates() error {
	// app.Log.Named("updates").Info("starting")
	qr, err := qbtCircuitBreaker.Execute(func() (*qbt.Response, error) {
		return app.Qbt.List()
	})
	if err != nil {
		if errors.Cause(err) == gobreaker.ErrOpenState {
			return nil
		}
		return errors.Wrap(err, "getting torrent list")
	}

	nzbs, err := app.Nzb.List()
	if err != nil {
		return errors.Wrap(err, "getting nzb list")
	}

	// check disk space every minute
	if time.Now().Unix()%60 == 0 {
		go checkDisk(nzbs.Status.FreeDiskSpaceMB, qr.AllPaused())
	}

	app.Events.Send("flame.qbittorrents", qr)
	app.Events.Send("flame.nzbs", nzbs)

	go updateQbittorrents(qr)
	go updateNzbs(nzbs)
	go func() {
		metrics := &Metrics{}
		metrics.Diskspace = fmt.Sprintf("%2.1f", float64(nzbs.Status.FreeDiskSpaceMB/1000))
		metrics.Torrents.DownloadRate = fmt.Sprintf("%2.1f", float64(qr.DownloadRate/1000))
		metrics.Torrents.UploadRate = fmt.Sprintf("%2.1f", float64(qr.UploadRate/1000))
		metrics.Nzbs.DownloadRate = fmt.Sprintf("%2.1f", float64(nzbs.Status.DownloadRate/1000))
		app.Events.Send("flame.metrics", metrics)
		app.Events.Send("flame.combined", &Combined{
			Torrents:  qr.Torrents,
			Nzbs:      nzbs.Result,
			NzbStatus: nzbs.Status,
			Metrics:   metrics,
		})
	}()

	return nil
}

func updateQbittorrents(resp *qbt.Response) {
	err := app.Cache.Set("flame-qbittorrents", resp)
	if err != nil {
		app.Workers.Log.Errorf("sendqbts: set cache failed: %s", err)
	}
}

func updateNzbs(resp *nzbget.GroupResponse) {
	err := app.Cache.Set("flame-nzbs", resp)
	if err != nil {
		app.Workers.Log.Errorf("sendnzbs: set cache failed: %s", err)
	}
}
