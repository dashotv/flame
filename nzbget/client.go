package nzbget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	URL string
}

func NewClient(endpoint string) *Client {
	client := &Client{
		URL: endpoint,
	}
	return client
}

type Response struct {
	APIVersion string `json:"version"`
}

func (c *Client) request(path string, params url.Values, target interface{}) (err error) {
	var url string
	var client *http.Client
	var request *http.Request
	var response *http.Response
	var body []byte

	url = fmt.Sprintf("%s/%s", c.URL, path)

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return errors.Wrap(err, "creating "+url+" request failed")
	}

	client = &http.Client{}
	if response, err = client.Do(request); err != nil {
		//log.Fatal(err)
		return errors.Wrap(err, "error making http request")
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		//log.Fatal(err)
		return errors.Wrap(err, "reading request body")
	}

	logrus.Debugf("body: %s", string(body))

	if target == nil {
		return nil
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	return nil
}
