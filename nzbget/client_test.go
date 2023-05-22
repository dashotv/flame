package nzbget

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/poll"
)

var nzbgetURL string

var name = "John Wick 4"
var link = "https://api.nzbgeek.info/api?t=get&id=842168bd21a2d26a86c5d04421fc8812&apikey=eISG7JzxXnmWskK632mjY3CHRylfVuiX"

func init() {
	nzbgetURL = os.Getenv("NZBGET_URL")
}

func TestClient_Add(t *testing.T) {
	c := NewClient(nzbgetURL)
	o := NewOptions()
	o.NiceName = name
	i, err := c.Add(link, o)
	require.NoError(t, err)
	require.GreaterOrEqual(t, i, int64(1))
}

func TestClient_List(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.Groups()
	require.NoError(t, err)
	require.NotNil(t, r)
	printList(r)
}

func TestClient_Pause(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.Groups()
	require.NoError(t, err)
	require.NotNil(t, r)

	err = c.Pause(r[0].ID)
	require.NoError(t, err)

	check := func(l poll.LogT) poll.Result {
		r, err = c.Groups()
		if err != nil {
			poll.Error(errors.Wrap(err, "failed to get groups"))
		}
		if r[0].Status == "PAUSED" || r[0].Status == "QUEUED" {
			return poll.Success()
		}
		return poll.Continue("Status != PAUSED (%s)", r[0].Status)
	}

	poll.WaitOn(t, check, poll.WithTimeout(30*time.Second), poll.WithDelay(1*time.Second))
}

func TestClient_Resume(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.Groups()
	require.NoError(t, err)
	require.NotNil(t, r)

	err = c.Resume(r[0].ID)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	r, err = c.Groups()
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "DOWNLOADING", r[0].Status)
	printList(r)
}

func TestClient_Remove(t *testing.T) {
	c := NewClient(nzbgetURL)
	r, err := c.Groups()
	require.NoError(t, err)
	require.NotNil(t, r)
	err = c.Remove(r[0].ID)
	require.NoError(t, err)
}

func TestClient_History(t *testing.T) {
	c := NewClient(nzbgetURL)
	list, err := c.History(false)
	require.NoError(t, err, "should not return error")
	printHistory(list)
}

func TestClient_Delete(t *testing.T) {
	c := NewClient(nzbgetURL)
	r, err := c.History(false)
	require.NoError(t, err)
	require.NotNil(t, r)
	err = c.Delete(r[0].ID)
	require.NoError(t, err)
	r, err = c.History(true)
	require.NoError(t, err)
	printHistory(r)
}

func TestClient_Destroy(t *testing.T) {
	c := NewClient(nzbgetURL)
	r, err := c.History(true)
	require.NoError(t, err)
	require.NotNil(t, r)
	printHistory(r)
	err = c.Destroy(r[0].ID)
	require.NoError(t, err)
	r, err = c.History(true)
	require.NoError(t, err)
	printHistory(r)
}

func TestClient_Status(t *testing.T) {
	c := NewClient(nzbgetURL)
	status, err := c.Status()
	require.NoError(t, err, "should not return error")
	fmt.Printf("status.DownloadedSizeMB: %d\n", status.DownloadedSizeMB)
}

func TestClient_Version(t *testing.T) {
	c := NewClient(nzbgetURL)
	v, err := c.Version()
	if err != nil {
		require.NoError(t, err, "should not return error")
	}
	require.Equal(t, "21.1", v)
}

func TestClient_PauseAll(t *testing.T) {
	c := NewClient(nzbgetURL)
	err := c.PauseAll()
	require.NoError(t, err, "should not return error")
}

func TestClient_ResumeAll(t *testing.T) {
	c := NewClient(nzbgetURL)
	err := c.ResumeAll()
	require.NoError(t, err, "should not return error")
}
