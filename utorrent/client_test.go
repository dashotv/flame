package utorrent

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	flameURL = os.Getenv("FLAME_URL")
	URLs     = map[string]string{
		"6f8cd699135b491513e65d967a052a7087750d9c": "http://www.slackware.com/torrents/slackware-14.1-install-d1.torrent",
		//"6f8cd699135b491513e65d967a052a7087750d9c": "./test/slackware-14.1-install-d1.torrent",
		"6293ed83c38a1b988f201f7f0b1ab315971ff38f": "magnet:?xt=urn:btih:6293ed83c38a1b988f201f7f0b1ab315971ff38f&dn=ubuntu+mini+remix+13.04+x32+&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80&tr=udp%3a%2f%2ftracker.publicbt.com%3a80&tr=udp%3a%2f%2ftracker.istole.it%3a6969&tr=udp%3a%2f%2ftracker.ccc.de%3a80&tr=udp%3a%2f%2fopen.demonii.com%3a1337",
		"6e5b9cb5635a95436963da0eab1f17a19e150168": "https://nyaa.si/download/920527.torrent",
		"684a8cb51634e4b82634c08be933ff32c44d8457": "https://nyaa.si/download/937262.torrent",
		//"b3c94d183bb461166419c1fe5111981d54ee84b0": "https://www.shanaproject.com/download/151055/",
	}

	magnetURLs = map[string]string{
		"6293ed83c38a1b988f201f7f0b1ab315971ff38f": "magnet:?xt=urn:btih:6293ed83c38a1b988f201f7f0b1ab315971ff38f&dn=ubuntu+mini+remix+13.04+x32+&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80&tr=udp%3a%2f%2ftracker.publicbt.com%3a80&tr=udp%3a%2f%2ftracker.istole.it%3a6969&tr=udp%3a%2f%2ftracker.ccc.de%3a80&tr=udp%3a%2f%2fopen.demonii.com%3a1337",
	}
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func TestClient_List(t *testing.T) {
	var r *Response
	var err error

	c := NewClient(flameURL)
	if r, err = c.List(); err != nil {
		require.NoError(t, err, "should be able to list")
	}
	printResponse(r)
}

func printResponse(r *Response) {
	fmt.Printf("DownloadRate: %4.3f UploadRate: %4.3f\n", r.DownloadRate, r.UploadRate)

	for _, t := range r.Torrents {
		fmt.Printf("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
		for _, f := range t.Files {
			fmt.Printf("  %d %6.2f%% %10.2fmb           %s\n", f.Priority, f.DownloadedPercent(), f.SizeMb(), f.Name)
		}
	}
}

func TestClient_Add(t *testing.T) {
	var r *Response
	var err error

	c := NewClient(flameURL)

	if r, err = c.List(); err != nil {
		require.NoError(t, err, "should be able to list")
	}
	printResponse(r)

	for k, v := range URLs {
		s, err := c.Add(v)
		require.NoError(t, err, "should be able to add: ", v)
		require.Equal(t, k, s, "hashes should match")

		time.Sleep(1 * time.Second)
		err = c.Remove(k, true)
		require.NoError(t, err, "shouldn't fail to remove")
	}

	time.Sleep(1 * time.Second)
	if r, err = c.List(); err != nil {
		require.NoError(t, err, "should be able to list")
	}
	printResponse(r)
}
