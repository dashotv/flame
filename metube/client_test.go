package metube

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var metubeURL string

func init() {
	godotenv.Load("../.env")
	metubeURL = os.Getenv("METUBE_URL")
}

func TestClient_Add(t *testing.T) {
	url := "https://www.dailymotion.com/embed/video/k4LoixovMHT0X5ycZJk?pubtool=oembed"
	c := New(metubeURL, true)
	err := c.Add("test", url)
	require.NoError(t, err)
}

func TestClient_History(t *testing.T) {
	c := New(metubeURL, true)
	res, err := c.History()
	require.NoError(t, err)
	require.NotNil(t, res)
	spew.Dump(res)
}
