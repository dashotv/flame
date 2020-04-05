package nzbget

import (
	"os"
)

var nzbgetURL string

func init() {
	nzbgetURL = os.Getenv("NZBGET_URL")
}
