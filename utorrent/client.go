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

type params map[string]interface{}

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
		return nil, c.error(err, "getting torrent list")
	}
	//fmt.Printf("parsed: %#v\n", parsed)
	r.Load(parsed)

	fileParams := url.Values{}
	fileParams.Add("action", "getfiles")
	for _, t := range r.Torrents {
		fileParams.Add("hash", t.Hash)
	}
	if err := c.request("", fileParams, files); err != nil {
		return nil, c.error(err, "getting torrent files")
	}
	r.LoadFiles(files)

	return r, nil
}

func (c *Client) Add(URI string) (string, error) {
	u, err := url.Parse(URI)
	if err != nil {
		return "", c.error(err, "could not add uri")
	}

	switch u.Scheme {
	case "http", "https":
		return c.addURL(u)
	case "magnet":
		return c.addMagnet(u)
	default:
		return "", c.error(nil, "only URL and magnet supported")
	}
}

func (c *Client) addMagnet(magnet *url.URL) (string, error) {
	i, err := magnetInfohash(magnet)
	if err != nil {
		return "", err
	}

	if err := c.action(params{"action": "add-url", "s": magnet.String()}); err != nil {
		return "", c.error(err, "could not add magnet")
	}

	return i, nil
}

func (c *Client) addFile(file string) (string, error) {
	i, err := fileInfohash(file)
	if err != nil {
		return "", err
	}

	return i, nil
}

func (c *Client) addURL(URL *url.URL) (string, error) {
	tmp, err := downloadURL(URL.String())
	if err != nil {
		return "", err
	}

	i, err := fileInfohash(tmp)
	if err != nil {
		return "", err
	}

	if err := c.action(params{"action": "add-url", "s": URL.String()}); err != nil {
		return "", c.error(err, "could not add magnet")
	}

	return i, nil
}

func (c *Client) Remove(infohash string, delete bool) error {
	a := "remove"
	if delete {
		a = "removedata"
	}

	if err := c.action(params{"action": a, "hash": infohash}); err != nil {
		return c.error(err, fmt.Sprintf("could not remove %s", infohash))
	}

	return nil
}

func (c *Client) Wanted(infohash string) (bool, error) {
	t, err := c.Get(infohash)
	if err != nil {
		return false, err
	}

	for _, f := range t.Files {
		if f.Priority != 0 {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) Get(infohash string) (*Torrent, error) {
	r, err := c.List()
	if err != nil {
		return nil, err
	}
	return r.Get(infohash), nil
}

func (c *Client) Want(infohash string, ids []int) error {
	return c.Priorities(infohash, 2, ids)
}

func (c *Client) WantNone(infohash string) error {
	t, err := c.Get(infohash)
	if err != nil {
		return err
	}

	ids := make([]int, len(t.Files))
	for i, f := range t.Files {
		ids[i] = f.Number
	}

	return c.Priorities(infohash, 0, ids)
}

func (c *Client) Priorities(infohash string, priority int, ids []int) error {
	return c.action(params{"action": "setprio", "p": priority, "f": ids, "hash": infohash})
}

func (c *Client) Pause(infohash string) error {
	return c.action(params{"action": "pause", "hash": infohash})
}

func (c *Client) PauseAll() error {
	return c.action(params{"action": "pauseall"})
}

func (c *Client) Resume(infohash string) error {
	return c.action(params{"action": "unpause", "hash": infohash})
}

func (c *Client) ResumeAll() error {
	return c.action(params{"action": "unpauseall"})
}

func (c *Client) Stop(infohash string) error {
	return c.action(params{"action": "stop", "hash": infohash})
}

func (c *Client) Start(infohash string) error {
	return c.action(params{"action": "start", "hash": infohash})
}

func (c *Client) Label(infohash, label string) error {
	return c.action(params{"action": "setprops", "hash": infohash, "s": "label", "v": label})
}

// Private

func (c *Client) authenticate() (err error) {
	var response *http.Response
	var doc *goquery.Document

	if c.authenticated {
		return nil
	}

	if response, err = http.Get(c.Url + "/token.html"); err != nil {
		return c.error(err, "getting token")
	}
	defer response.Body.Close()

	// get token from response
	if doc, err = goquery.NewDocumentFromResponse(response); err != nil {
		return c.error(err, "reading http response")
	}
	if c.token = doc.Find("div#token").Text(); c.token == "" {
		return c.error(nil, "token not found")
	}

	// find GUID cookie and store value
	for _, cookie := range response.Cookies() {
		if cookie.Name == "GUID" {
			c.cookie = cookie.Value
			c.authenticated = true
			return nil
		}
	}

	return c.error(nil, "failed to authenticate")
}

//func (c *Client) action(params url.Values) error {
func (c *Client) action(params params) error {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return c.request("", values, nil)
}

func (c *Client) request(action string, params url.Values, target map[string]interface{}) (err error) {
	var url string
	var client *http.Client
	var request *http.Request
	var response *http.Response
	var body []byte

	if err = c.authenticate(); err != nil {
		return c.error(err, "authentication failed")
	}

	url = fmt.Sprintf("%s/%s", c.Url, action)

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return c.error(err, "creating "+url+" request failed")
	}

	request.Header.Set("Cookie", fmt.Sprintf("GUID=%s; count=1", c.cookie))
	params.Set("cid", "1")
	params.Set("token", c.token)
	request.URL.RawQuery = params.Encode()
	logrus.Debugf("request: %s", request.URL.String())

	client = &http.Client{}
	if response, err = client.Do(request); err != nil {
		//log.Fatal(err)
		return c.error(err, "error making http request")
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		//log.Fatal(err)
		return c.error(err, "reading request body")
	}

	logrus.Debugf("body: %s", string(body))

	if target == nil {
		return nil
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return c.error(err, "json unmarshall")
	}

	return nil
}

func (c *Client) error(err error, msg string) error {
	c.token = ""

	if err != nil {
		return errors.Wrap(err, msg)
	}

	return errors.New(msg)
}
