package nzbget

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var nzbgetURL string

func init() {
	nzbgetURL = os.Getenv("NZBGET_URL")
}

func TestClient_Add(t *testing.T) {
	c := NewClient(nzbgetURL)
	i, err := c.Add("https://api.nzbgeek.info/api?t=get&apikey=2b2f7303f77672ad619df8589e88b5d3&id=c4d0c7dee6bab6a35d9b74592ade0bb7")
	require.NoError(t, err)
	require.GreaterOrEqual(t, i, int64(1))
}

func TestClient_List(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.List()
	require.NoError(t, err)
	require.NotNil(t, r)
	printList(r)
}

func TestClient_Pause(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.List()
	require.NoError(t, err)
	require.NotNil(t, r)

	err = c.Pause(r[0].ID)
	require.NoError(t, err)

	time.Sleep(10 * time.Second)
	r, err = c.List()
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "PAUSED", r[0].Status)
	printList(r)
}

func TestClient_Resume(t *testing.T) {
	c := NewClient(nzbgetURL)

	r, err := c.List()
	require.NoError(t, err)
	require.NotNil(t, r)

	err = c.Resume(r[0].ID)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	r, err = c.List()
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "DOWNLOADING", r[0].Status)
	printList(r)
}

func TestClient_Remove(t *testing.T) {
	c := NewClient(nzbgetURL)
	r, err := c.List()
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
	require.Equal(t, "21.0", v)
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
