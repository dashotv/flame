package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
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

func (s *Server) Updates() {
	qbt, err := qb.List()
	if err != nil {
		s.Log.Errorf("couldn't get torrent list: %s", err)
		return
	}

	go s.updateQbittorrents(qbt)

	nzbs, err := nzb.List()
	if err != nil {
		s.Log.Errorf("couldn't get nzb list: %s", err)
		return
	}

	go s.updateNzbs(nzbs)
	go s.checkDisk(nzbs)

	go func() {
		metrics := &Metrics{}
		metrics.Diskspace = fmt.Sprintf("%2.1f", float64(nzbs.Status.FreeDiskSpaceMB/1000))
		metrics.Torrents.DownloadRate = fmt.Sprintf("%2.1f", float64(qbt.DownloadRate/1000))
		metrics.Torrents.UploadRate = fmt.Sprintf("%2.1f", float64(qbt.UploadRate/1000))
		metrics.Nzbs.DownloadRate = fmt.Sprintf("%2.1f", float64(nzbs.Status.DownloadRate/1000))
		s.metricsChannel <- metrics
		s.combined <- &Combined{
			Torrents: qbt.Torrents,
			Nzbs:     nzbs.Result,
			Metrics:  metrics,
		}
	}()
}

func (s *Server) updateQbittorrents(resp *qbt.Response) {
	b, err := json.Marshal(&resp)
	if err != nil {
		s.Log.Errorf("couldn't marshal torrents: %s", err)
		return
	}

	status := cache.Set(ctx, "flame-qbittorrents", string(b), time.Minute)
	if status.Err() != nil {
		s.Log.Errorf("sendqbts: set cache failed: %s", status.Err())
	}
	s.qbtChannel <- resp
}

func (s *Server) updateNzbs(resp *nzbget.GroupResponse) {
	b, err := json.Marshal(&resp)
	if err != nil {
		s.Log.Errorf("couldn't marshal nzbs: %s", err)
		return
	}

	status := cache.Set(ctx, "flame-nzbs", string(b), time.Minute)
	if status.Err() != nil {
		s.Log.Errorf("sendnzbs: set cache failed: %s", status.Err())
	}
	s.nzbChannel <- resp
}
