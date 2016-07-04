package flame

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/grengojbo/goquery"
)

type Client struct {
	Url string

	token         string
	cookie        string
	authenticated bool
}

func (c *Client) Connect(url *string) {
	c.Url = *url
	c.token = ""
	c.cookie = ""
	c.authenticated = false
}

func (c *Client) Authenticate() error {
	if !c.authenticated {
		response, err := http.Get(c.Url + "/token.html")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// get token from response
		doc, err := goquery.NewDocumentFromResponse(response)
		if err != nil {
			return err
		}
		c.token = doc.Find("div#token").Text()

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

func (c *Client) List(response *Response) error {
	parsed := make(map[string]interface{})

	err := c.get(&parsed)
	if err != nil {
		return err
	}
	//fmt.Printf("parsed: %q\n", parsed)
	response.Load(&parsed)

	return nil
}

func (c *Client) get(target *map[string]interface{}) error {
	if !c.authenticated {
		err := c.Authenticate()
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s/?token=%s&cid=1&list=1", c.Url, c.token)
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	request.Header.Set("Cookie", fmt.Sprintf("GUID=%s; count=1", c.cookie))
	response, err := client.Do(request)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	//fmt.Printf("body: %s\n", string(body))
	//fmt.Print("\n")

	json.Unmarshal(body, target)

	return nil
}
