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
	c := Client{}
	r := Response{}
	c.Connect(&url)
	c.List(&r)

	for _, t := range r.Torrents {
		str := fmt.Sprintf("%3.0f %.2f %s\n", t.Queue, t.Progress, t.Name)
		fmt.Println(str)
		//publish(str)
	}
}
