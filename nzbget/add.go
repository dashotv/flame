package nzbget

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/andrewstuart/go-nzb"
	"github.com/pkg/errors"
)

func (c *Client) Add(URL string) (int64, error) {
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

	r, err := c.rpc.Call("append", name, enc, "", 0, false, false, "", 0, "SCORE", []int{})
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

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "could not read file")
	}

	return string(b), nil
}

func base64encode(s string) string {
	data := []byte(s)
	str := base64.StdEncoding.EncodeToString(data)
	return str
}

func downloadURL(URL string) (string, error) {
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		return "", errors.Wrap(err, "could not http get url")
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("/tmp", "flame-download-*")
	if err != nil {
		return "", errors.Wrap(err, "could not get tmp file")
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "could not copy file")
	}

	return file.Name(), nil
}

func nzbName(data string) (string, error) {
	nzb := nzb.NZB{}
	err := xml.Unmarshal([]byte(data), &nzb)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal")
	}

	return nzb.Meta["name"], nil
}
