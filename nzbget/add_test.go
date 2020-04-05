package nzbget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Add(t *testing.T) {
	c := NewClient(nzbgetURL)
	i, err := c.Add("https://api.nzbgeek.info/api?t=get&apikey=2b2f7303f77672ad619df8589e88b5d3&id=c4d0c7dee6bab6a35d9b74592ade0bb7")
	require.NoError(t, err)
	require.GreaterOrEqual(t, i, int64(1))
}
