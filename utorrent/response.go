package utorrent

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Build        float64
	Torrents     []*Torrent
	UploadRate   float64
	DownloadRate float64
	CacheId      string
	Timestamp    time.Time
}

func NewResponse() *Response {
	return &Response{
		Timestamp: time.Now(),
	}
}

func (r *Response) Get(hash string) *Torrent {
	for _, t := range r.Torrents {
		if t.Hash == hash {
			return t
		}
	}
	return nil
}

func (r *Response) Load(value map[string]interface{}) {
	r.Build = value["build"].(float64)

	if val, ok := value["torrentc"]; ok {
		r.CacheId = val.(string)
	}

	if val, ok := value["torrents"]; ok {
		for _, t := range val.([]interface{}) {
			//fmt.Println(t)
			torrent := &Torrent{}
			torrent.Load(t.([]interface{}))
			r.UploadRate += torrent.UploadRate / 1024
			r.DownloadRate += torrent.DownloadRate / 1024
			r.Torrents = append(r.Torrents, torrent)
		}
	}
}

func (r *Response) LoadFiles(value map[string]interface{}) {
	if val, ok := value["files"]; ok {
		list := val.([]interface{})
		for i := 0; i+1 <= len(list); i = i + 2 {
			logrus.Debugf("%s = %#v", list[i], list[i+1])
			for _, f := range list[i+1].([]interface{}) {
				file := &File{}
				file.Load(f.([]interface{}))
				t := r.Get(list[i].(string))
				t.AddFile(file)
			}
		}
	}
}

func (r *Response) Count() int {
	return len(r.Torrents)
}

func (r *Response) JSON() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
