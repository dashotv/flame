package nzbget

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func printHistory(list []History) {
	for _, s := range list {
		fmt.Printf("%5d %5d %s\n", s.ID, s.NZBID, s.Name)
	}
}

func TestClient_History(t *testing.T) {
	c := NewClient(nzbgetURL)
	list, err := c.History()
	if err != nil {
		require.NoError(t, err, "should not return error")
	}
	printHistory(list)
}
