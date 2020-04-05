package utorrent

import (
	"net/url"
	"testing"
)

func Test_magnetInfohash(t *testing.T) {
	for k, v := range magnetURLs {
		u, err := url.Parse(v)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		i, err := magnetInfohash(u)
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if i != k {
			t.Errorf("infohash: wanted %s, got %s", k, i)
		}
	}
}
