package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/minion"
)

type Metrics struct {
	Diskspace string `json:"diskspace"`
	Torrents  struct {
		DownloadRate string `json:"download_rate"`
		UploadRate   string `json:"upload_rate"`
	} `json:"torrents"`
	Nzbs struct {
		DownloadRate string `json:"download_rate"`
	} `json:"nzbs"`
}

type Combined struct {
	Torrents  []*qbt.Torrent
	Nzbs      []nzbget.Group
	NzbStatus nzbget.Status
	Metrics   *Metrics
}

type Updates struct {
	minion.WorkerDefaults[*Updates]
}

func (j *Updates) Kind() string { return "updates" }
func (j *Updates) Work(ctx context.Context, job *minion.Job[*Updates]) error {
	// app.Log.Named("updates").Info("starting")
	qbt, err := app.Qbt.List()
	if err != nil {
		return errors.Wrap(err, "getting torrent list")
	}

	go j.updateQbittorrents(qbt)

	nzbs, err := app.Nzb.List()
	if err != nil {
		return errors.Wrap(err, "getting nzb list")
	}

	go j.updateNzbs(nzbs)
	go j.checkDisk(nzbs, qbt)

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

func (j *Updates) updateQbittorrents(resp *qbt.Response) {
	b, err := json.Marshal(&resp)
	if err != nil {
		app.Workers.Log.Errorf("couldn't marshal torrents: %s", err)
		return
	}

	err = app.Cache.Set("flame-qbittorrents", string(b))
	if err != nil {
		app.Workers.Log.Errorf("sendqbts: set cache failed: %s", err)
	}

	app.Events.Send("flame.qbittorrents", resp)
}

func (j *Updates) updateNzbs(resp *nzbget.GroupResponse) {
	b, err := json.Marshal(&resp)
	if err != nil {
		app.Workers.Log.Errorf("couldn't marshal nzbs: %s", err)
		return
	}

	err = app.Cache.Set("flame-nzbs", string(b))
	if err != nil {
		app.Workers.Log.Errorf("sendnzbs: set cache failed: %s", err)
	}
	app.Events.Send("flame.nzbs", resp)
}

// pauseAll pauses all torrents and sets a flag in the cache
func (j *Updates) pauseAll() error {
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

func (j *Updates) resumeAll() error {
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
func (j *Updates) diskPaused() bool {
	paused := "false"
	_, err := app.Cache.Get("flame-disk-paused", &paused)
	if err != nil {
		// app.Workers.Log.Errorf("paused: %s", err)
		return false
	}
	return paused == "true"
}

// allPaused checks if all torrents are paused and the paused flag is true
func (j *Updates) allPaused() bool {
	if !j.diskPaused() {
		return false
	}

	ok, err := app.Qbt.AllPaused()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to check if all qbts are paused: %s", err)
	}
	return ok
}

func (j *Updates) checkDisk(resp *nzbget.GroupResponse, qbt *qbt.Response) {
	if resp.Status.FreeDiskSpaceMB < 25000 && !j.allPaused() {
		err := j.pauseAll()
		if err != nil {
			app.Workers.Log.Errorf("checkdisk: failed to pause all qbts: %s", err)
		}
		return
	}

	if !j.diskPaused() {
		return
	}
	if !qbt.AllPaused() {
		return
	}

	err := j.resumeAll()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}
}
