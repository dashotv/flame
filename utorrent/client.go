package utorrent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/grengojbo/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Url string

	token         string
	cookie        string
	authenticated bool
}

func NewClient(url string) *Client {
	return &Client{
		Url:           url,
		token:         "",
		cookie:        "",
		authenticated: false,
	}
}

func (c *Client) List() (*Response, error) {
	r := NewResponse()
	parsed := make(map[string]interface{})
	files := make(map[string]interface{})

	params := url.Values{}
	params.Add("list", "1")
	if err := c.request("", params, parsed); err != nil {
		return nil, errors.Wrap(err, "getting torrent list")
	}
	//fmt.Printf("parsed: %#v\n", parsed)
	r.Load(parsed)

	fileParams := url.Values{}
	fileParams.Add("action", "getfiles")
	for _, t := range r.Torrents {
		fileParams.Add("hash", t.Hash)
	}
	if err := c.request("", fileParams, files); err != nil {
		return nil, errors.Wrap(err, "getting torrent files")
	}
	r.LoadFiles(files)

	return r, nil
}

// Private

func (c *Client) authenticate() (err error) {
	var response *http.Response
	var doc *goquery.Document

	if c.authenticated {
		return nil
	}

	if response, err = http.Get(c.Url + "/token.html"); err != nil {
		return errors.Wrap(err, "getting token")
	}
	defer response.Body.Close()

	// get token from response
	if doc, err = goquery.NewDocumentFromResponse(response); err != nil {
		return errors.Wrap(err, "reading http response")
	}
	if c.token = doc.Find("div#token").Text(); c.token == "" {
		return errors.New("token not found")
	}

	// find GUID cookie and store value
	for _, cookie := range response.Cookies() {
		if cookie.Name == "GUID" {
			c.cookie = cookie.Value
			c.authenticated = true
			return nil
		}
	}

	return errors.New("failed to authenticate")
}

func (c *Client) request(action string, params url.Values, target map[string]interface{}) (err error) {
	var url string
	var client *http.Client
	var request *http.Request
	var response *http.Response
	var body []byte

	if err = c.authenticate(); err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	url = fmt.Sprintf("%s/%s", c.Url, action)

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return errors.Wrap(err, "creating "+url+" request failed")
	}

	request.Header.Set("Cookie", fmt.Sprintf("GUID=%s; count=1", c.cookie))
	params.Set("cid", "1")
	params.Set("token", c.token)
	request.URL.RawQuery = params.Encode()
	logrus.Debugf("request: %s", request.URL.String())

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

	if err = json.Unmarshal(body, &target); err != nil {
		return errors.Wrap(err, "json unmarshall")
	}

	return nil
}
