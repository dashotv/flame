package flame

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	flameURL = os.Getenv("FLAME_URL")
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func TestClient_List(t *testing.T) {
	var r *Response
	var err error

	c := NewClient(flameURL)
	if r, err = c.List(); err != nil {
		t.Error(err)
	}

	fmt.Printf("DownloadRate: %4.3f UploadRate: %4.3f\n", r.DownloadRate, r.UploadRate)

	for _, t := range r.Torrents {
		fmt.Printf("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
		for _, f := range t.Files {
			fmt.Printf("  %d %6.2f%% %10.2fmb           %s\n", f.Priority, f.DownloadedPercent(), f.SizeMb(), f.Name)
		}
	}
}
