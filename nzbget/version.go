package nzbget

import (
	"fmt"
	"net/url"
)

type VersionResponse struct {
	*Response
	Version string `json:"result"`
}

func (c *Client) Version() (string, error) {
	version := &VersionResponse{}
	err := c.request("version", url.Values{}, version)
	if err != nil {
		return "", err
	}
	fmt.Printf("response: %#v\n", version)
	return version.Version, nil
}
