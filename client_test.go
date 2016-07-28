package flame

import (
	"fmt"
	"testing"
	"os"
)

var (
	url = os.Getenv("FLAME_URL")
)

func TestClient_List(t *testing.T) {
	var r *Response
	var err error

	c := NewClient(url)
	if r, err = c.List(); err != nil {
		t.Error(err)
	}

	for _, t := range r.Torrents {
		fmt.Printf("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
	}
}
