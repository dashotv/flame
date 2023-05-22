package qbt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

type Api struct {
	Client        *ClientWithResponses
	URL           string
	Authenticated bool
	Ctx           context.Context
	Jar           *cookiejar.Jar
	Sid           string
}

type Response struct {
	Torrents     []*Torrent
	UploadRate   int // bytes / s
	DownloadRate int // bytes / s
	RequestID    int `json:"rid"`
	Timestamp    time.Time
}

func clientCookieJar(jar http.CookieJar) func(*Client) error {
	return func(c *Client) error {
		h := http.DefaultClient

		h.Jar = jar
		c.Client = h
		return nil
	}
}

func NewApi(url string) *Api {
	api := &Api{}

	// ensure url ends with "/"
	if url[len(url)-1:] != "/" {
		url += "/"
	}

	url += "api/v2/"

	api.URL = url

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	api.Jar = jar

	client, err := NewClientWithResponses(url, clientCookieJar(jar))
	if err != nil {
		panic(err)
	}
	api.Client = client
	api.Ctx = context.Background()

	return api
}

func (a *Api) Login(username, password string) (ok bool, err error) {
	params := &AuthLoginPostParams{}
	body := AuthLoginPostFormdataRequestBody{Username: username, Password: password}
	resp, err := a.Client.AuthLoginPostWithFormdataBody(a.Ctx, params, body)
	if err != nil {
		return false, err
	}

	for _, c := range resp.Cookies() {
		if c.Name == "SID" {
			a.Sid = c.Value
			a.Authenticated = true
			break
		}
	}

	return a.Authenticated, nil
}

func (a *Api) List() (*Response, error) {
	sync, err := a.Sync(0)

	torrents, err := a.Torrents("priority")
	if err != nil {
		return nil, err
	}

	out := &Response{
		Torrents:     torrents,
		RequestID:    sync.Rid,
		DownloadRate: sync.ServerState.DlInfoSpeed,
		UploadRate:   sync.ServerState.UpInfoSpeed,
		Timestamp:    time.Now(),
	}

	return out, nil
}

func (a *Api) Torrents(sort string) ([]*Torrent, error) {
	var torrents []*Torrent
	filter := TorrentsInfoPostRequestFilterAll

	req := TorrentsInfoPostFormdataRequestBody{
		Filter: &filter,
		Sort:   &sort,
	}
	resp, err := a.Client.TorrentsInfoPostWithFormdataBody(a.Ctx, req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&torrents)
	if err != nil {
		return nil, err
	}

	for _, t := range torrents {
		files, err := a.TorrentFiles(t.Hash)
		if err != nil {
			return nil, err
		}
		t.Files = files
	}

	return torrents, nil
}

func (a *Api) TorrentFiles(hash string) ([]*TorrentFile, error) {
	var files []*TorrentFile

	req := TorrentsFilesPostRequest{
		Hash: hash,
	}
	resp, err := a.Client.TorrentsFilesPostWithFormdataBody(a.Ctx, req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	for i, f := range files {
		f.ID = i
		f.Progress = f.Progress * 100
	}

	return files, nil
}

func (a *Api) Torrent(hash string) (*Torrent, error) {
	torrents, err := a.Torrents("priority")
	if err != nil {
		return nil, err
	}

	for _, t := range torrents {
		if t.Hash == hash {
			return t, nil
		}
	}

	return nil, fmt.Errorf("could not find torrent for infohash: %s", hash)
}

func (a *Api) Sync(requestID int) (*Sync, error) {
	s := &Sync{}
	rid := int64(requestID)

	body := SyncMaindataPostFormdataRequestBody{Rid: &rid}
	resp, err := a.Client.SyncMaindataPostWithFormdataBody(a.Ctx, body)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (a *Api) Add(link string, options map[string]string) (string, error) {
	hash, err := InfohashFromURL(link)
	if err != nil {
		return "", err
	}

	params := map[string]string{"urls": link}
	for k, v := range options {
		params[k] = v
	}

	buffer, content, err := setupParams(params)
	if err != nil {
		return "", errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsAddPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("add error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return hash, nil
}

func (a *Api) Delete(hash string, perm bool) error {
	params := map[string]string{"hashes": hash, "deleteFiles": fmt.Sprintf("%t", perm)}

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsDeletePostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete error status code %d: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (a *Api) Pause(hash string) error {
	params := map[string]string{"hashes": hash}

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsPausePostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pause error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) PauseAll() error {
	return a.Pause("all")
}

func (a *Api) Resume(hash string) error {
	params := map[string]string{"hashes": hash}

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsResumePostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resume error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) ResumeAll() error {
	return a.Resume("all")
}

func (a *Api) SetTag(hashes []string, tag string) error {
	params := processInfoHashList(hashes)
	params["tag"] = tag

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsAddTagsPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tag error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) SetLabel(hashes []string, label string) error {
	return a.SetTag(hashes, label)
}

func (a *Api) Recheck(hashes []string) error {
	params := processInfoHashList(hashes)

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsRecheckPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recheck error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) IncreasePriority(hashes []string) error {
	params := processInfoHashList(hashes)

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsIncreasePrioPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("increase priority error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) DecreasePriority(hashes []string) error {
	params := processInfoHashList(hashes)

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsDecreasePrioPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("decrease priority error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) Want(hash string, fileIDs []string) error {
	return a.SetFilePriority(hash, fileIDs, "1")
}

func (a *Api) WantNone(hash string) error {
	files, err := a.TorrentFiles(hash)
	if err != nil {
		return err
	}

	ids := make([]string, len(files))
	for _, f := range files {
		ids = append(ids, fmt.Sprintf("%d", f.ID))
	}

	return a.SetFilePriority(hash, ids, "0")
}

func (a *Api) Wanted(hash string) (bool, error) {
	files, err := a.TorrentFiles(hash)
	if err != nil {
		return false, err
	}

	for _, f := range files {
		if f.Priority > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (a *Api) WantedCount(hash string) (int, error) {
	count := 0
	files, err := a.TorrentFiles(hash)
	if err != nil {
		return -1, err
	}

	for _, f := range files {
		if f.Priority > 0 {
			count++
		}
	}

	return count, nil
}

func (a *Api) SetFilePriority(hash string, fileIDs []string, priority string) error {
	if !validPriority(priority) {
		return fmt.Errorf("bad priority %s", priority)
	}

	list := []string{}
	for _, s := range fileIDs {
		if s != "" {
			list = append(list, s)
		}
	}

	ids := strings.Join(list, "|")

	params := make(map[string]string)
	params["hash"] = hash
	params["id"] = ids
	params["priority"] = priority

	buffer, content, err := setupParams(params)
	if err != nil {
		return errors.Wrap(err, "setup params")
	}

	resp, err := a.Client.TorrentsFilePrioPostWithBody(a.Ctx, content, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set file priority error status code %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}
