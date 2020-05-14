// shamelessly stolen from https://github.com/superturkey650/go-qbittorrent
// Updated to support API v2.5.1 (Qbittorrent v4.2.5)
package qbt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"net/url"
	"strconv"
	"strings"

	wrapper "github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

//ErrBadPriority means the priority is not allowd by qbittorrent
var ErrBadPriority = errors.New("priority not available")

//ErrBadResponse means that qbittorrent sent back an unexpected response
var ErrBadResponse = errors.New("received bad response")

//Client creates a connection to qbittorrent and performs requests
type Client struct {
	http          *http.Client
	URL           string
	Authenticated bool
	Jar           http.CookieJar
}

//NewClient creates a new client connection to qbittorrent
func NewClient(url string) *Client {
	client := &Client{}

	// ensure url ends with "/"
	if url[len(url)-1:] != "/" {
		url += "/"
	}

	url += "api/v2/"

	client.URL = url

	// create cookie jar
	client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.http = &http.Client{
		Jar: client.Jar,
	}
	return client
}

//Login logs you in to the qbittorrent client
//returns the current authentication status
func (client *Client) Login(username string, password string) (loggedIn bool, err error) {
	credentials := make(map[string]string)
	credentials["username"] = username
	credentials["password"] = password

	resp, err := client.post("auth/login", credentials)
	if err != nil {
		return false, err
	} else if resp.Status != "200 OK" { // check for correct status code
		//b, _ := ioutil.ReadAll(resp.Body)
		//fmt.Printf("body: %s\n", string(b))
		return false, wrapper.Wrap(ErrBadResponse, "couldnt log in")
	}

	// change authentication status so we know were authenticated in later requests
	client.Authenticated = true

	// add the cookie to cookie jar to authenticate later requests
	if cookies := resp.Cookies(); len(cookies) > 0 {
		cookieURL, _ := url.Parse("http://localhost:8080")
		client.Jar.SetCookies(cookieURL, cookies)
	}

	// create a new client with the cookie jar and replace the old one
	// so that all our later requests are authenticated
	client.http = &http.Client{
		Jar: client.Jar,
	}

	return client.Authenticated, nil
}

//Logout logs you out of the qbittorrent client
//returns the current authentication status
func (client *Client) Logout() (loggedOut bool, err error) {
	resp, err := client.get("auth/logout", nil)
	if err != nil {
		return false, err
	}

	// check for correct status code
	if resp.Status != "200 OK" {
		return false, wrapper.Wrap(ErrBadResponse, "couldnt log out")
	}

	// change authentication status so we know were not authenticated in later requests
	client.Authenticated = false

	return client.Authenticated, nil
}

//Shutdown shuts down the qbittorrent client
func (client *Client) Shutdown() (shuttingDown bool, err error) {
	resp, err := client.get("command/shutdown", nil)

	// return true if successful
	return resp.Status == "200 OK", err
}

type Response struct {
	Torrents     []*Torrent
	UploadRate   int // bytes / s
	DownloadRate int // bytes / s
	RequestID    int `json:"rid"`
	Timestamp    time.Time
}

func (r *Response) Pretty() string {
	return fmt.Sprintf("%3d v %d b/s ^ %d b/s torrents:%d time:%s", r.RequestID, r.DownloadRate, r.UploadRate, len(r.Torrents), r.Timestamp)
}

func (client *Client) List() (*Response, error) {
	sync, err := client.Sync("0")
	if err != nil {
		return nil, err
	}

	torrents, err := client.Torrents(map[string]string{"sort": "priority"})
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Torrents:     torrents,
		RequestID:    sync.Rid,
		DownloadRate: sync.ServerState.DlInfoSpeed,
		UploadRate:   sync.ServerState.UpInfoSpeed,
		Timestamp:    time.Now(),
	}
	return resp, nil
}

//Torrents returns a list of all torrents in qbittorrent matching your filter
func (client *Client) Torrents(filters map[string]string) (torrentList []*Torrent, err error) {
	torrents := make([]*Torrent, 0)

	resp, err := client.get("torrents/info", filters)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&torrents)
	if err != nil {
		return nil, err
	}

	for _, t := range torrents {
		files, err := client.TorrentFiles(t.Hash)
		if err != nil {
			return nil, err
		}
		t.Files = files
	}

	return torrents, nil
}

//Torrent is a convenience function returning the matching torrent
func (client *Client) Torrent(infohash string) (*Torrent, error) {
	torrents, err := client.Torrents(nil)
	if err != nil {
		return nil, err
	}

	for _, t := range torrents {
		if t.Hash == infohash {
			return t, nil
		}
	}

	return nil, fmt.Errorf("could not find torrent for infohash: %s", infohash)
}

//Torrent returns a specific torrent matching the infoHash
//func (client *Client) Torrent(infoHash string) (*Torrent, error) {
//	t := &Torrent{}
//	resp, err := client.get("torrents/properties", map[string]string{"hash": infoHash})
//	if err != nil {
//		return t, err
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&t)
//	if err != nil {
//		return nil, err
//	}
//
//	return t, nil
//}

//TorrentTrackers returns all trackers for a specific torrent matching the infoHash
//func (client *Client) TorrentTrackers(infoHash string) ([]Tracker, error) {
//	var t []Tracker
//	resp, err := client.get("torrents/trackers", map[string]string{"hash": infoHash})
//	if err != nil {
//		return t, err
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&t)
//	if err != nil {
//		return nil, err
//	}
//
//	return t, nil
//}

//TorrentWebSeeds returns seeders for a specific torrent matching the infoHash
//func (client *Client) TorrentWebSeeds(infoHash string) ([]WebSeed, error) {
//	var w []WebSeed
//	resp, err := client.get("torrents/webseeds", map[string]string{"hash": infoHash})
//	if err != nil {
//		return w, err
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&w)
//	if err != nil {
//		return nil, err
//	}
//
//	return w, nil
//}

//TorrentFiles gets the files of a specifc torrent matching the infoHash
func (client *Client) TorrentFiles(infoHash string) ([]*TorrentFile, error) {
	files := make([]*TorrentFile, 0)
	resp, err := client.get("torrents/files", map[string]string{"hash": infoHash})
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	for i, f := range files {
		f.ID = i
	}

	return files, nil
}

//Sync returns the server state and list of torrents in one struct
func (client *Client) Sync(rid string) (Sync, error) {
	var s Sync

	params := make(map[string]string)
	params["rid"] = rid

	resp, err := client.get("sync/maindata", params)
	if err != nil {
		return s, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (client *Client) Add(link string, options map[string]string) (string, error) {
	infohash, err := InfohashFromURL(link)
	if err != nil {
		return "", err
	}

	_, err = client.DownloadFromLink(link, options)
	if err != nil {
		return "", err
	}

	return infohash, nil
}

//DownloadFromLink starts downloading a torrent from a link
func (client *Client) DownloadFromLink(link string, options map[string]string) (*http.Response, error) {
	options["urls"] = link
	return client.postMultipartData("torrents/add", options)
}

//DownloadFromFile starts downloading a torrent from a file
func (client *Client) DownloadFromFile(file string, options map[string]string) (*http.Response, error) {
	return client.postMultipartFile("torrents/add", file, options)
}

//AddTrackers adds trackers to a specific torrent matching infoHash
func (client *Client) AddTrackers(infoHash string, trackers string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	params["urls"] = trackers

	return client.post("command/addTrackers", params)
}

//Pause pauses a specific torrent matching infoHash
func (client *Client) Pause(infoHash string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)

	return client.post("command/pause", params)
}

//PauseAll pauses all torrents
func (client *Client) PauseAll() (*http.Response, error) {
	return client.get("command/pause", map[string]string{"hashes": "all"})
}

////PauseMultiple pauses a list of torrents matching the infoHashes
//func (client *Client) PauseMultiple(infoHashList []string) (*http.Response, error) {
//	params := client.processInfoHashList(infoHashList)
//	return client.post("command/pauseAll", params)
//}

//Resume resumes a specific torrent matching infoHash
func (client *Client) Resume(infoHash string) (*http.Response, error) {
	params := make(map[string]string)
	params["hash"] = strings.ToLower(infoHash)
	return client.post("command/resume", params)
}

//ResumeAll resumes all torrents matching infoHashes
func (client *Client) ResumeAll() (*http.Response, error) {
	return client.get("command/resume", map[string]string{"hashes": "all"})
}

//ResumeMultiple resumes a list of torrents matching infoHashes
//func (client *Client) ResumeMultiple(infoHashList []string) (*http.Response, error) {
//	params := client.processInfoHashList(infoHashList)
//	return client.post("command/resumeAll", params)
//}

//SetLabel is an alias for SetTag
func (client *Client) SetLabel(infoHashList []string, label string) (*http.Response, error) {
	return client.SetTag(infoHashList, label)
}

//SetTag sets a tag for a list of torrents matching infoHashes
func (client *Client) SetTag(infoHashList []string, tag string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["tags"] = tag

	return client.post("torrents/addTags", params)
}

//SetCategory sets the category for a list of torrents matching infoHashes
//func (client *Client) SetCategory(infoHashList []string, category string) (*http.Response, error) {
//	params := client.processInfoHashList(infoHashList)
//	params["category"] = category
//
//	return client.post("command/setLabel", params)
//}

func (client *Client) Delete(infohash string, perm bool) (*http.Response, error) {
	if perm {
		return client.DeletePermanently([]string{infohash})
	}
	return client.DeleteTemp([]string{infohash})
}

//DeleteTemp deletes the temporary files for a list of torrents matching infoHashes
func (client *Client) DeleteTemp(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.get("torrents/delete", params)
}

//DeletePermanently deletes all files for a list of torrents matching infoHashes
func (client *Client) DeletePermanently(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["deleteFiles"] = "true"
	return client.get("torrents/delete", params)
}

//Recheck rechecks a list of torrents matching infoHashes
func (client *Client) Recheck(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.post("torrents/recheck", params)
}

//IncreasePriority increases the priority of a list of torrents matching infoHashes
func (client *Client) IncreasePriority(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.post("torrents/increasePrio", params)
}

//DecreasePriority decreases the priority of a list of torrents matching infoHashes
func (client *Client) DecreasePriority(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.post("torrents/decreasePrio", params)
}

//SetMaxPriority sets the max priority for a list of torrents matching infoHashes
func (client *Client) SetMaxPriority(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.post("torrents/topPrio", params)
}

//SetMinPriority sets the min priority for a list of torrents matching infoHashes
func (client *Client) SetMinPriority(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.post("torrents/bottomPrio", params)
}

//Want is a convenience function for setting file priorities for given torrent
func (client *Client) Want(infohash string, fileIDs []string) error {
	_, err := client.SetFilePriority(infohash, fileIDs, "1")
	return err
}

//WantNone is a convenience function for setting file priorities for given torrent
func (client *Client) WantNone(infohash string) error {
	files, err := client.TorrentFiles(infohash)
	if err != nil {
		return err
	}

	ids := make([]string, len(files))
	for _, f := range files {
		ids = append(ids, fmt.Sprintf("%d", f.ID))
	}

	_, err = client.SetFilePriority(infohash, ids, "0")
	return err
}

//Wanted is a convenience function that returns true if any file of given torrent is wanted
func (client *Client) Wanted(infohash string) (bool, error) {
	torrent, err := client.Torrent(infohash)
	if err != nil {
		return false, err
	}

	for _, f := range torrent.Files {
		if f.Priority > 0 {
			return true, nil
		}
	}

	return false, nil
}

//SetFilePriority sets the file priorities for torrent matching infoHash
func (client *Client) SetFilePriority(infoHash string, fileIDs []string, priority string) (*http.Response, error) {
	// disallow certain priorities that are not allowed by the WEBUI API
	priorities := [...]string{"0", "1", "2", "7"}
	for _, v := range priorities {
		if v == priority {
			return nil, ErrBadPriority
		}
	}

	ids := strings.Join(fileIDs, "|")

	params := make(map[string]string)
	params["hash"] = infoHash
	params["id"] = ids
	params["priority"] = priority

	return client.post("torrents/filePrio", params)
}

//GetTorrentDownloadLimit gets the download limit for a list of torrents
func (client *Client) GetTorrentDownloadLimit(infoHashList []string) (limits map[string]string, err error) {
	var l map[string]string
	params := client.processInfoHashList(infoHashList)
	resp, err := client.post("torrents/downloadLimit", params)
	if err != nil {
		return l, err
	}
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return l, err
	}
	return l, nil
}

//SetTorrentDownloadLimit sets the download limit for a list of torrents
func (client *Client) SetTorrentDownloadLimit(infoHashList []string, limit string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["limit"] = limit
	return client.post("torrents/setDownloadLimit", params)
}

//SetTorrentShareLimit sets the share ratio and seeding time limits for a list of torrents
func (client *Client) SetTorrentShareLimit(infoHashList []string, ratioLimit string, seedingTime string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["ratioLimit"] = ratioLimit
	params["seedingTimeLimit"] = seedingTime
	return client.post("torrents/setShareLimits", params)
}

//GetTorrentUploadLimit gets the upload limit for a list of torrents
func (client *Client) GetTorrentUploadLimit(infoHashList []string) (limits map[string]string, err error) {
	var l map[string]string
	params := client.processInfoHashList(infoHashList)
	resp, err := client.post("torrents/uploadLimit", params)
	if err != nil {
		return l, err
	}
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return l, err
	}
	return l, nil
}

//SetTorrentUploadLimit sets the upload limit of a list of torrents
func (client *Client) SetTorrentUploadLimit(infoHashList []string, limit string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["limit"] = limit
	return client.post("torrents/setUploadLimit", params)
}

//SetPreferences sets the preferences of your qbittorrent client
//func (client *Client) SetPreferences(params map[string]string) (*http.Response, error) {
//	return client.post("command/setPreferences", params)
//}

//ToggleSequentialDownload toggles the download sequence of a list of torrents
func (client *Client) ToggleSequentialDownload(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.get("torrents/toggleSequentialDownload", params)
}

//ToggleFirstLastPiecePriority toggles first last piece priority of a list of torrents
func (client *Client) ToggleFirstLastPiecePriority(infoHashList []string) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	return client.get("torrents/toggleFirstLastPiecePrio", params)
}

//ForceStart force starts a list of torrents
func (client *Client) ForceStart(infoHashList []string, value bool) (*http.Response, error) {
	params := client.processInfoHashList(infoHashList)
	params["value"] = strconv.FormatBool(value)
	return client.post("torrents/setForceStart", params)
}

//GetGlobalDownloadLimit gets the global download limit of your qbittorrent client
func (client *Client) GetGlobalDownloadLimit() (limit int, err error) {
	var l int
	resp, err := client.get("transfer/downloadLimit", nil)
	if err != nil {
		return l, err
	}
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return l, err
	}
	return l, nil
}

//SetGlobalDownloadLimit sets the global download limit of your qbittorrent client
func (client *Client) SetGlobalDownloadLimit(limit string) (*http.Response, error) {
	params := make(map[string]string)
	params["limit"] = limit
	return client.post("transfer/setDownloadLimit", params)
}

//GetGlobalUploadLimit gets the global upload limit of your qbittorrent client
func (client *Client) GetGlobalUploadLimit() (limit int, err error) {
	var l int
	resp, err := client.get("transfer/uploadLimit", nil)
	if err != nil {
		return l, err
	}
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return l, err
	}
	return l, nil
}

//SetGlobalUploadLimit sets the global upload limit of your qbittorrent client
func (client *Client) SetGlobalUploadLimit(limit string) (*http.Response, error) {
	params := make(map[string]string)
	params["limit"] = limit
	return client.post("transfer/setUploadLimit", params)
}

//GetAlternativeSpeedStatus gets the alternative speed status of your qbittorrent client
func (client *Client) GetAlternativeSpeedStatus() (status bool, err error) {
	var s int
	resp, err := client.get("transfer/speedLimitsMode", nil)
	if err != nil {
		return false, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return false, err
	}

	return s == 1, nil
}

//ToggleAlternativeSpeed toggles the alternative speed of your qbittorrent client
func (client *Client) ToggleAlternativeSpeed() (*http.Response, error) {
	return client.get("transfer/toggleSpeedLimitsMode", nil)
}
