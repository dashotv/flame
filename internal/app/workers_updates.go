package app

import (
	"fmt"

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

	go updateQbittorrents(qbt)

	nzbs, err := app.Nzb.List()
	if err != nil {
		return errors.Wrap(err, "getting nzb list")
	}

	go updateNzbs(nzbs)
	go checkDisk(nzbs, qbt)

	go func() {
		metrics := &Metrics{}
		metrics.Diskspace = fmt.Sprintf("%2.1f", float64(nzbs.Status.FreeDiskSpaceMB/1000))
		metrics.Torrents.DownloadRate = fmt.Sprintf("%2.1f", float64(qbt.DownloadRate/1000))
		metrics.Torrents.UploadRate = fmt.Sprintf("%2.1f", float64(qbt.UploadRate/1000))
		metrics.Nzbs.DownloadRate = fmt.Sprintf("%2.1f", float64(nzbs.Status.DownloadRate/1000))
		app.Events.Send("flame.metrics", metrics)
		app.Events.Send("flame.combined", &Combined{
			Torrents: qbt.Torrents,
			Nzbs:     nzbs.Result,
			Metrics:  metrics,
		})
	}()

	return nil
}

func updateQbittorrents(resp *qbt.Response) {
	err := app.Cache.Set("flame-qbittorrents", resp)
	if err != nil {
		app.Workers.Log.Errorf("sendqbts: set cache failed: %s", err)
	}

	app.Events.Send("flame.qbittorrents", resp)
}

func updateNzbs(resp *nzbget.GroupResponse) {
	err := app.Cache.Set("flame-nzbs", resp)
	if err != nil {
		app.Workers.Log.Errorf("sendnzbs: set cache failed: %s", err)
	}
	app.Events.Send("flame.nzbs", resp)
}

// pauseAll pauses all torrents and sets a flag in the cache
func pauseAll() error {
	err := app.Qbt.PauseAll()
	if err != nil {
		return errors.Wrap(err, "pausing all")
	}

	err = app.Cache.Set("flame-disk-paused", true)
	if err != nil {
		app.Workers.Log.Errorf("sendqbts: set cache failed: %s", err)
	}

	return nil
}

func resumeAll() error {
	err := app.Qbt.ResumeAll()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}

	err = app.Cache.Delete("flame-disk-paused")
	if err != nil {
		app.Workers.Log.Errorf("sendqbts: set cache failed: %s", err)
	}

	return nil
}

// diskPaused checks the cache for the flag
func diskPaused() bool {
	paused := "false"
	_, err := app.Cache.Get("flame-disk-paused", &paused)
	if err != nil {
		// app.Workers.Log.Errorf("paused: %s", err)
		return false
	}
	return paused == "true"
}

// allPaused checks if all torrents are paused and the paused flag is true
func allPaused() bool {
	if !diskPaused() {
		return false
	}

	ok, err := app.Qbt.AllPaused()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to check if all qbts are paused: %s", err)
	}
	return ok
}

func checkDisk(resp *nzbget.GroupResponse, qbt *qbt.Response) {
	if resp.Status.FreeDiskSpaceMB < 25000 && !allPaused() {
		err := pauseAll()
		if err != nil {
			app.Workers.Log.Errorf("checkdisk: failed to pause all qbts: %s", err)
		}
		return
	}

	if !diskPaused() {
		return
	}
	if !qbt.AllPaused() {
		return
	}

	err := resumeAll()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}
}
