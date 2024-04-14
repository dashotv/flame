package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

func Updates() error {
	// app.Log.Named("updates").Info("starting")
	qbt, err := app.Qbt.List()
	if err != nil {
		return errors.Wrap(err, "getting torrent list")
	}

	nzbs, err := app.Nzb.List()
	if err != nil {
		return errors.Wrap(err, "getting nzb list")
	}

	// check disk space every minute
	if time.Now().Unix()%60 == 0 {
		go checkDisk(nzbs.Status.FreeDiskSpaceMB, qbt.AllPaused())
	}

	app.Events.Send("flame.qbittorrents", qbt)
	app.Events.Send("flame.nzbs", nzbs)

	go updateQbittorrents(qbt)
	go updateNzbs(nzbs)
	go func() {
		metrics := &Metrics{}
		metrics.Diskspace = fmt.Sprintf("%2.1f", float64(nzbs.Status.FreeDiskSpaceMB/1000))
		metrics.Torrents.DownloadRate = fmt.Sprintf("%2.1f", float64(qbt.DownloadRate/1000))
		metrics.Torrents.UploadRate = fmt.Sprintf("%2.1f", float64(qbt.UploadRate/1000))
		metrics.Nzbs.DownloadRate = fmt.Sprintf("%2.1f", float64(nzbs.Status.DownloadRate/1000))
		app.Events.Send("flame.metrics", metrics)
		app.Events.Send("flame.combined", &Combined{
			Torrents:  qbt.Torrents,
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
