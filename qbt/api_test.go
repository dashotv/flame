package qbt

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../.env")
}

var (
	URLs = map[string]string{
		"6f8cd699135b491513e65d967a052a7087750d9c": "http://www.slackware.com/torrents/slackware-14.1-install-d1.torrent",
		//"6f8cd699135b491513e65d967a052a7087750d9c": "./test/slackware-14.1-install-d1.torrent",
		"6293ed83c38a1b988f201f7f0b1ab315971ff38f": "magnet:?xt=urn:btih:6293ed83c38a1b988f201f7f0b1ab315971ff38f&dn=ubuntu+mini+remix+13.04+x32+&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80&tr=udp%3a%2f%2ftracker.publicbt.com%3a80&tr=udp%3a%2f%2ftracker.istole.it%3a6969&tr=udp%3a%2f%2ftracker.ccc.de%3a80&tr=udp%3a%2f%2fopen.demonii.com%3a1337",
		"6e5b9cb5635a95436963da0eab1f17a19e150168": "https://nyaa.si/download/920527.torrent",
		"684a8cb51634e4b82634c08be933ff32c44d8457": "https://nyaa.si/download/937262.torrent",
		//"b3c94d183bb461166419c1fe5111981d54ee84b0": "https://www.shanaproject.com/download/151055/",
	}

	magnetURLs = map[string]string{
		"6293ed83c38a1b988f201f7f0b1ab315971ff38f": "magnet:?xt=urn:btih:6293ed83c38a1b988f201f7f0b1ab315971ff38f&dn=ubuntu+mini+remix+13.04+x32+&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80&tr=udp%3a%2f%2ftracker.publicbt.com%3a80&tr=udp%3a%2f%2ftracker.istole.it%3a6969&tr=udp%3a%2f%2ftracker.ccc.de%3a80&tr=udp%3a%2f%2fopen.demonii.com%3a1337",
		"68d583f28dc697401a8b769fecbe011272e1dc89": "magnet:?xt=urn:btih:68d583f28dc697401a8b769fecbe011272e1dc89",
	}
)

func TestApi_Login(t *testing.T) {
	url := os.Getenv("QBITTORRENT_URL")
	pw := os.Getenv("QBITTORRENT_PASSWORD")
	assert.NotEmpty(t, url)
	assert.NotEmpty(t, pw)

	api := NewApi(url)
	ok, err := api.Login("admin", pw)
	assert.True(t, ok)
	assert.NoError(t, err)

	r, err := api.List()
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestApi_Torrents(t *testing.T) {
	url := os.Getenv("QBITTORRENT_URL")
	assert.NotEmpty(t, url)

	api := NewApi(url)
	assert.NotNil(t, api)

	list, err := api.Torrents("priority")
	assert.NoError(t, err)
	assert.NotNil(t, list)
}

func TestApi_Add(t *testing.T) {
	url := os.Getenv("QBITTORRENT_URL")
	assert.NotEmpty(t, url)

	api := NewApi(url)
	opts := map[string]string{}

	list, err := api.Torrents("priority")
	assert.NoError(t, err)
	count := len(list)

	for k, v := range URLs {
		s, err := api.Add(v, opts)
		assert.NoError(t, err, "should be able to add: ", v)
		assert.Equal(t, k, s, "hashes should match")

		time.Sleep(1 * time.Second)
		err = api.Delete(k, true)
		assert.NoError(t, err, "shouldn't fail to remove")
	}

	for k, v := range magnetURLs {
		s, err := api.Add(v, opts)
		assert.NoError(t, err, "should be able to add: ", v)
		assert.Equal(t, k, s, "hashes should match")

		time.Sleep(1 * time.Second)
		err = api.Delete(k, true)
		assert.NoError(t, err, "shouldn't fail to remove")
	}

	list, err = api.Torrents("priority")
	assert.NoError(t, err)
	assert.Equal(t, len(list), count)
}
