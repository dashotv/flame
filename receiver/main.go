package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame/utorrent"
	"github.com/dashotv/mercury"
)

type TorrentsByIndex []*utorrent.Torrent

func (a TorrentsByIndex) Len() int           { return len(a) }
func (a TorrentsByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TorrentsByIndex) Less(i, j int) bool { return a[i].Queue < a[j].Queue }

func main() {
	m, err := mercury.New("flame", nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting receiver...")
	channel := make(chan *utorrent.Response, 5)
	m.Receiver("flame.torrents", channel)
	channel2 := make(chan string, 5)
	m.Receiver("flame.downloads", channel2)

	for {
		select {
		case r := <-channel:
			//logrus.Infof("received message")
			sort.Sort(TorrentsByIndex(r.Torrents))
			for _, t := range r.Torrents {
				logrus.Infof("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
				for _, f := range t.Files {
					logrus.Infof("  %d %6.2f%% %10.2fmb           %s\n", f.Priority, f.DownloadedPercent(), f.SizeMb(), f.Name)
				}
			}
		case s := <-channel2:
			fmt.Printf("received2: %#v\n", s)
		case <-time.After(30 * time.Second):
			fmt.Println("timeout")
			os.Exit(0)
		}
	}
}
