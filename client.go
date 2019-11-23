package flame

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grengojbo/goquery"
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
	r := &Response{}
	parsed := make(map[string]interface{})

	if err := c.request("", "list=1", &parsed); err != nil {
		return nil, err
	}

	//fmt.Println("parsed: ", parsed)
	r.Load(&parsed)

	return r, nil
}

// Private

func (c *Client) authenticate() (err error) {
	var response *http.Response
	var doc *goquery.Document

	if !c.authenticated {
		if response, err = http.Get(c.Url + "/token.html"); err != nil {
			return err
		}
		defer response.Body.Close()

		// get token from response
		if doc, err = goquery.NewDocumentFromResponse(response); err != nil {
			return err
		}
		if c.token = doc.Find("div#token").Text(); c.token == "" {
			return fmt.Errorf("token not found")
		}

		// find GUID cookie and store value
		for _, cookie := range response.Cookies() {
			if cookie.Name == "GUID" {
				c.cookie = cookie.Value
				c.authenticated = true
				return nil
			}
		}
		return fmt.Errorf("failed to authenticate")
	} else {
		return nil
	}
}

func (c *Client) request(action string, params string, target *map[string]interface{}) (err error) {
	var url string
	var client *http.Client
	var request *http.Request
	var response *http.Response
	var body []byte

	if err = c.authenticate(); err != nil {
		return err
	}

	url = fmt.Sprintf("%s/%s?%s&cid=1&token=%s", c.Url, action, params, c.token)
	//fmt.Printf("request: %s\n", url)

	client = &http.Client{}

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		//log.Fatal(err)
		return err
	}

	request.Header.Set("Cookie", fmt.Sprintf("GUID=%s; count=1", c.cookie))

	if response, err = client.Do(request); err != nil {
		//log.Fatal(err)
		return err
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		//log.Fatal(err)
		return err
	}

	//fmt.Println("body: ", string(body))

	json.Unmarshal(body, target)

	return nil
}
