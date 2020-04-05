package nzbget

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Status(t *testing.T) {
	c := NewClient(nzbgetURL)
	status, err := c.Status()
	if err != nil {
		require.NoError(t, err, "should not return error")
	}
	fmt.Printf("status.DownloadedSizeMB: %d\n", status.DownloadedSizeMB)
}
