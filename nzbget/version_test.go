package nzbget

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Version(t *testing.T) {
	c := NewClient(nzbgetURL)
	v, err := c.Version()
	if err != nil {
		require.NoError(t, err, "should not return error")
	}

	fmt.Printf("version: %s\n", v)
}
