package qbt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var qb *Client
var user, pass string
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
	}
)

func init() {
	qb = NewClient("http://qbittorrent.dasho.net/")
	user = "admin"
	pass = "adminadmin"
}

func printList(list *Response) {
	fmt.Println(list.Pretty())
	for _, t := range list.Torrents {
		fmt.Println(t.Pretty())
	}
}

func TestClient_Login(t *testing.T) {
	ok, err := qb.Login(user, pass)
	assert.NoError(t, err, "logging in")
	assert.Equal(t, true, ok, "logged in")
}

//func TestClient_Logout(t *testing.T) {
//	ok, err := qb.Logout()
//	if err != nil {
//		t.Errorf("error logging out: %s", err)
//	}
//	if ok {
//		t.Error("logged out but still logged in")
//	}
//}

func TestClient_Add(t *testing.T) {
	//if r, err = c.List(); err != nil {
	//	require.NoError(t, err, "should be able to list")
	//}
	//printResponse(r)

	opts := map[string]string{}

	for k, v := range URLs {
		s, err := qb.Add(v, opts)
		assert.NoError(t, err, "should be able to add: ", v)
		assert.Equal(t, k, s, "hashes should match")

		time.Sleep(1 * time.Second)
		_, err = qb.Delete(k, true)
		assert.NoError(t, err, "shouldn't fail to remove")
	}

	//time.Sleep(1 * time.Second)
	//if r, err = c.List(); err != nil {
	//	require.NoError(t, err, "should be able to list")
	//}
	//printResponse(r)
}


func TestClient_WantNone(t *testing.T) {
	//if r, err = c.List(); err != nil {
	//	require.NoError(t, err, "should be able to list")
	//}
	//printResponse(r)

	opts := map[string]string{}
	hash := "6f8cd699135b491513e65d967a052a7087750d9c"
	url := URLs[hash]

	s, err := qb.Add(url, opts)
	assert.NoError(t, err, "should be able to add: ", hash)
	assert.Equal(t, hash, s, "hashes should match")

	time.Sleep(1 * time.Second)
	TestClient_List(t)
	err = qb.WantNone(hash)
	assert.NoError(t, err, "should be able to want none: ", hash)

	time.Sleep(1 * time.Second)
	torrent, err := qb.Torrent(hash)
	for i, file := range torrent.Files {
		if file.Priority != 0 {
			assert.NoError(t, err, "should not be wanted: ", i)
		}
	}

	time.Sleep(1 * time.Second)
	TestClient_List(t)
	_, err = qb.Delete(hash, true)
	assert.NoError(t, err, "shouldn't fail to remove")
}

func TestClient_Sync(t *testing.T) {
	sync, err := qb.Sync("0")
	assert.NoError(t, err, "should be able to get sync")
	assert.Equal(t, 1, sync.Rid, "request id")
}

func TestClient_Torrents(t *testing.T) {
	list, err := qb.Torrents(nil)
	require.NoError(t, err, "should be able to get torrents")
	fmt.Printf("list: %#v\n", list)
}

func TestClient_List(t *testing.T) {
	list, err := qb.List()
	require.NoError(t, err, "should be able to get list")
	printList(list)
}
