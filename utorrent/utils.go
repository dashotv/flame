package utorrent

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
)

func query(params map[string]string) string {
	list := []string{}
	for k, v := range params {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(list, "&")
}

func getString(value interface{}) string {
	if value != nil {
		return value.(string)
	}
	return ""
}

func getFloat64(value interface{}) float64 {
	if value != nil {
		return value.(float64)
	}
	return 0.0
}

func getInt(value interface{}) int {
	if value != nil {
		return value.(int)
	}
	return 0
}

func downloadURL(URL string) (string, error) {
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("/tmp", "flame-download-*")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	return file.Name(), err
}

type torrentFile struct {
	Info struct {
		Name        string `bencode:"name"`
		Length      int64  `bencode:"length"`
		MD5Sum      string `bencode:"md5sum,omitempty"`
		PieceLength int64  `bencode:"piece length"`
		Pieces      string `bencode:"pieces"`
		Private     bool   `bencode:"private,omitempty"`
	} `bencode:"info"`

	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list,omitempty"`
	CreationDate int64       `bencode:"creation date,omitempty"`
	Comment      string      `bencode:"comment,omitempty"`
	CreatedBy    string      `bencode:"created by,omitempty"`
	URLList      interface{} `bencode:"url-list,omitempty"`
}

func fileInfohash(path string) (string, error) {
	mi, err := metainfo.LoadFromFile(path)
	if err != nil {
		return "", errors.Wrap(err, "could not read metainfo")
	}

	return mi.HashInfoBytes().String(), nil
}

func magnetInfohash(URL *url.URL) (string, error) {
	q, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return "", errors.Wrap(err, "could not parse query")
	}

	h := ""
	for _, v := range q["xt"] {
		s := strings.Split(v, ":")
		if s[0] == "urn" && s[1] == "btih" {
			h = s[2]
			break
		}
	}

	return h, nil
}
