package nzbget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame/internal/jsonrpc"
)

const PriorityVeryLow = -100
const PriorityLow = -50
const PriorityNormal = 0
const PriorityHigh = 50
const PriorityVeryHigh = 100
const PriorityForce = 900

type Client struct {
	URL string
	rpc jsonrpc.RPCClient
}

type AppendOptions struct {
	NiceName   string
	Category   string
	Priority   int
	AddToTop   bool
	AddPaused  bool
	DupeKey    string
	DupeScore  int
	DupeMode   string
	Parameters []struct {
		Name  string
		Value string
	}
}

func NewOptions() *AppendOptions {
	return &AppendOptions{
		Category: "",
		Priority: PriorityNormal,
		DupeMode: "SCORE",
	}
}

func NewClient(endpoint string) *Client {
	client := &Client{
		URL: endpoint,
		rpc: jsonrpc.NewClient(endpoint),
	}
	return client
}

type Response struct {
	APIVersion string `json:"version"`
	Error      string
	Status     *Status
	Timestamp  time.Time
}

func (c *Client) List() (*GroupResponse, error) {
	s, err := c.Status()
	if err != nil {
		return nil, err
	}

	r := &GroupResponse{Response: &Response{Timestamp: time.Now(), Status: s}}
	err = c.request("listgroups", nil, r)
	return r, err
}

func (c *Client) Groups() ([]Group, error) {
	r, err := c.List()
	if err != nil {
		return nil, errors.Wrap(err, "could not get list")
	}
	return r.Result, nil
}

func (c *Client) Remove(number int) error {
	// group delete
	return c.EditQueue("GroupDelete", "", []int{number})
}

func (c *Client) Delete(number int) error {
	return c.EditQueue("HistoryDelete", "", []int{number})
}

func (c *Client) Destroy(number int) error {
	return c.EditQueue("HistoryFinalDelete", "", []int{number})
}

func (c *Client) Pause(number int) error {
	return c.EditQueue("GroupPause", "", []int{number})
}

func (c *Client) Resume(number int) error {
	return c.EditQueue("GroupResume", "", []int{number})
}

func (c *Client) PauseAll() error {
	r, err := c.rpc.Call("pausedownload", nil)
	if err != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Error != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Result != true {
		return errors.New("response result is not true")
	}
	return nil
}

func (c *Client) ResumeAll() error {
	r, err := c.rpc.Call("resumedownload", nil)
	if err != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Error != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Result != true {
		return errors.New("response result is not true")
	}
	return nil
}

func (c *Client) EditQueue(command, param string, ids []int) error {
	r, err := c.rpc.Call("editqueue", command, param, ids)
	if err != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Error != nil {
		return errors.Wrap(err, "could not pause all")
	}
	if r.Result != true {
		return errors.New("response result is not true")
	}
	return nil
}

func (c *Client) Add(URL string, options *AppendOptions) (int64, error) {
	path, err := downloadURL(URL)
	if err != nil {
		return 0, errors.Wrap(err, "could not download url")
	}

	str, err := readFile(path)
	if err != nil {
		return 0, errors.Wrap(err, "could not read downloaded file")
	}

	name, err := nzbName(str)
	if err != nil {
		return 0, errors.Wrap(err, "could not get nzb name")
	}
	enc := base64encode(str)

	if options.NiceName != "" {
		name = options.NiceName
	}

	r, err := c.rpc.Call("append", name, enc, options.Category, options.Priority, options.AddToTop, options.AddPaused, options.DupeKey, options.DupeScore, options.DupeMode, options.Parameters)
	if err != nil {
		if r != nil && r.Error != nil {
			return 0, errors.Wrap(err, r.Error.Error())
		}
		return 0, err
	}

	n := r.Result.(json.Number)
	i, err := n.Int64()
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (c *Client) History(hidden bool) ([]History, error) {
	r := &HistoryResponse{}
	err := c.request("history", url.Values{"": []string{fmt.Sprintf("%t", hidden)}}, r)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}

func (c *Client) Status() (*Status, error) {
	r := &StatusResponse{}
	err := c.request("status", url.Values{}, r)
	if err != nil {
		return r.Result, err
	}
	return r.Result, nil
}

func (c *Client) Version() (string, error) {
	version := &VersionResponse{}
	err := c.request("version", url.Values{}, version)
	if err != nil {
		return "", err
	}
	return version.Version, nil
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
	request.URL.RawQuery = params.Encode()

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
