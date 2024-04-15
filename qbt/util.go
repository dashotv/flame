package qbt

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
)

func validPriority(priority string) bool {
	switch priority {
	case "0", "1", "2", "7":
		return true
	default:
		return false
	}
}

func InfohashFromURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	switch u.Scheme {
	case "magnet":
		return magnetInfohash(u)
	case "http", "https":
		return httpInfohash(u)
	}

	return "", nil
}

func downloadURL(URL *url.URL) (string, error) {
	// Get the data
	resp, err := http.Get(URL.String())
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

func fileInfohash(path string) (string, error) {
	mi, err := metainfo.LoadFromFile(path)
	if err != nil {
		return "", errors.Wrap(err, "could not read metainfo")
	}

	return mi.HashInfoBytes().String(), nil
}

func httpInfohash(URL *url.URL) (string, error) {
	f, err := downloadURL(URL)
	if err != nil {
		return "", err
	}

	return fileInfohash(f)
}

func magnetInfohash(URL *url.URL) (string, error) {
	magnet, err := metainfo.ParseMagnetUri(URL.String())
	if err != nil {
		return "", err
	}

	// 	q, err := url.ParseQuery(URL.RawQuery)
	// 	if err != nil {
	// 		return "", errors.Wrap(err, "could not parse query")
	// 	}
	//
	// 	h := ""
	// 	for _, v := range q["xt"] {
	// 		s := strings.Split(v, ":")
	// 		if s[0] == "urn" && s[1] == "btih" {
	// 			h = s[2]
	// 			break
	// 		}
	// 	}

	return magnet.InfoHash.HexString(), nil
}

// processInfoHashList puts list into a combined (single element) map with all hashes connected with '|'
// this is how the WEBUI API recognizes multiple hashes
func processInfoHashList(infoHashList []string) (hashMap map[string]string) {
	params := map[string]string{}
	params["hashes"] = strings.Join(infoHashList, "|")
	return params
}

func setupParams(params map[string]string) (*bytes.Buffer, string, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	for key, val := range params {
		writer.WriteField(key, val)
	}
	if err := writer.Close(); err != nil {
		return nil, "", errors.Wrap(err, "failed to close writer")
	}

	return buffer, writer.FormDataContentType(), nil
}
