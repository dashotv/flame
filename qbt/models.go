package qbt

import (
	"encoding/json"
	"fmt"
)

// Torrent holds a basic torrent object from qbittorrent (from torrents/list)
type Torrent struct {
	AddedOn           int     `json:"added_on"`    // unixtime
	AmountLeft        int     `json:"amount_left"` // bytes
	AutoManaged       bool    `json:"auto_tmm"`
	Availability      float64 `json:"availability"` // percentage
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"` // unixtime
	DownloadLimit     int     `json:"dl_limit"`      // bytes/s
	DownloadSpeed     int     `json:"dlspeed"`       // bytes/s
	Downloaded        int     `json:"downloaded"`    // bytes
	Eta               int     `json:"eta"`           // seconds
	FirstLastPriority bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int     `json:"last_activity"` // unixtime
	MagnetURI         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"` // seconds
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"` // percentage / 100
	Ratio             float64 `json:"ratio"`
	RatioLimit        float64 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int     `json:"size"` // bytes
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`        // comma delimited
	TimeActive        int     `json:"time_active"` // seconds
	TotalSize         int     `json:"total_size"`  // bytes
	Tracker           string  `json:"tracker"`
	Uploaded          int     `json:"uploaded"` // bytes
	UploadLimit       int     `json:"up_limit"` // bytes / s
	UploadSpeed       int     `json:"upspeed"`  // bytes / s
	//SeenComplete time.Time `json:"seen_complete"`
	//SeedingTimeLimit int `json:"seeding_time_limit"`
	//SessionDownloaded int     `json:"downloaded_session"`
	//SessionUploaded int `json:"uploaded_session"` // bytes

	Files []*TorrentFile
}

type TorrentJSON struct {
	Hash         string
	Status       int
	State        string
	Name         string
	Size         float64 // in bytes
	Progress     float64 // float64 in per mils
	Downloaded   float64 // in bytes
	Uploaded     float64 // in bytes
	Ratio        float64 // float64 in per mils
	UploadRate   float64 // float64 in bytes / second
	DownloadRate float64 // float64 in bytes / second
	Finish       float64 // float64 seconds
	Label        string
	//PeersConnected float64
	//PeersTotal     float64
	//SeedsConnected float64
	//SeedsTotal     float64
	//Availability   float64 // in 1/65535ths
	Queue float64
	//Remaining      float64 // in bytes
	Files []*TorrentFile
}

func (t *Torrent) Pretty() string {
	return fmt.Sprintf("%3d %6.2f%% %12.2fmb %12.12s %s\n%s", t.Priority, t.Progress*100, t.SizeMb(), t.State, t.Name, t.filesPretty())
}
func (t *Torrent) filesPretty() string {
	s := ""
	for _, f := range t.Files {
		s += "  " + f.Pretty()
	}
	return s
}

func (t *Torrent) SizeMb() float64 {
	return float64(t.TotalSize) / 1000.0
}

func (t *Torrent) MarshalJSON() ([]byte, error) {
	out := &TorrentJSON{}
	out.Hash = t.Hash
	out.State = t.State
	out.Name = t.Name
	out.Size = float64(t.Size)
	out.Progress = t.Progress * 100
	out.Downloaded = float64(t.Downloaded)
	out.Uploaded = float64(t.Uploaded)
	out.Ratio = t.Ratio
	out.UploadRate = float64(t.UploadSpeed)
	out.DownloadRate = float64(t.DownloadSpeed)
	out.Finish = float64(t.Eta)
	out.Label = t.Category
	out.Queue = float64(t.Priority)
	out.Files = t.Files
	return json.Marshal(out)
}

//
////Torrent holds a torrent object from qbittorrent
////with more information than BasicTorrent
//type Torrent struct {
//	AdditionDate           int     `json:"addition_date"`
//	Comment                string  `json:"comment"`
//	CompletionDate         int     `json:"completion_date"`
//	CreatedBy              string  `json:"created_by"`
//	CreationDate           int     `json:"creation_date"`
//	DlLimit                int     `json:"dl_limit"`
//	DlSpeed                int     `json:"dl_speed"`
//	DlSpeedAvg             int     `json:"dl_speed_avg"`
//	Eta                    int     `json:"eta"`
//	LastSeen               int     `json:"last_seen"`
//	NbConnections          int     `json:"nb_connections"`
//	NbConnectionsLimit     int     `json:"nb_connections_limit"`
//	Peers                  int     `json:"peers"`
//	PeersTotal             int     `json:"peers_total"`
//	PieceSize              int     `json:"piece_size"`
//	PiecesHave             int     `json:"pieces_have"`
//	PiecesNum              int     `json:"pieces_num"`
//	Reannounce             int     `json:"reannounce"`
//	SavePath               string  `json:"save_path"`
//	SeedingTime            int     `json:"seeding_time"`
//	Seeds                  int     `json:"seeds"`
//	SeedsTotal             int     `json:"seeds_total"`
//	ShareRatio             float64 `json:"share_ratio"`
//	TimeElapsed            int     `json:"time_elapsed"`
//	TotalDownloaded        int     `json:"total_downloaded"`
//	TotalDownloadedSession int     `json:"total_downloaded_session"`
//	TotalSize              int     `json:"total_size"`
//	TotalUploaded          int     `json:"total_uploaded"`
//	TotalUploadedSession   int     `json:"total_uploaded_session"`
//	TotalWasted            int     `json:"total_wasted"`
//	UpLimit                int     `json:"up_limit"`
//	UpSpeed                int     `json:"up_speed"`
//	UpSpeedAvg             int     `json:"up_speed_avg"`
//}

//Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg      string `json:"msg"`
	NumPeers int    `json:"num_peers"`
	Status   string `json:"status"`
	URL      string `json:"url"`
}

//WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

//TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Size     int     `json:"size"`     // bytes
	Progress float64 `json:"progress"` // percent / 100
	Priority int     `json:"priority"`
	IsSeed   bool    `json:"is_seed"`
}

func (f *TorrentFile) Pretty() string {
	return fmt.Sprintf("%3d %6.2f%% %s\n", f.Priority, f.Progress*100, f.Name)
}

//Sync holds the sync response struct which contains
//the server state and a map of infohashes to Torrents
type Sync struct {
	Rid        int                `json:"rid"`
	FullUpdate bool               `json:"full_update"`
	Torrents   map[string]Torrent `json:"torrents"`
	//Categories  []string           `json:"categories"`
	ServerState struct {
		ConnectionStatus  string `json:"connection_status"`
		DhtNodes          int    `json:"dht_nodes"`
		DlInfoData        int    `json:"dl_info_data"`
		DlInfoSpeed       int    `json:"dl_info_speed"`
		DlRateLimit       int    `json:"dl_rate_limit"`
		Queueing          bool   `json:"queueing"`
		RefreshInterval   int    `json:"refresh_interval"`
		UpInfoData        int    `json:"up_info_data"`
		UpInfoSpeed       int    `json:"up_info_speed"`
		UpRateLimit       int    `json:"up_rate_limit"`
		UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	} `json:"server_state"`
}
