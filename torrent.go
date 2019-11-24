package flame

//http://help.utorrent.com/customer/en/portal/articles/1573947-torrent-labels-list---webapi

type Torrent struct {
	Hash           string
	Status         int
	State          string
	Name           string
	Size           float64 // in bytes
	Progress       float64 // float64eger in per mils
	Downloaded     float64 // in bytes
	Uploaded       float64 // in bytes
	Ratio          float64 // float64eger in per mils
	UploadRate     float64 // float64eger in bytes / second
	DownloadRate   float64 // float64eger in bytes / second
	Finish         float64 // float64eger seconds
	Label          string
	PeersConnected float64
	PeersTotal     float64
	SeedsConnected float64
	SeedsTotal     float64
	Availability   float64 // in 1/65535ths
	Queue          float64
	Remaining      float64 // in bytes
}

func (t *Torrent) Load(values []interface{}) {
	t.Hash = values[0].(string)
	t.Status = int(values[1].(float64))
	t.Name = values[2].(string)
	t.Size = values[3].(float64)
	t.Progress = values[4].(float64) / 10
	t.Downloaded = values[5].(float64)
	t.Uploaded = values[6].(float64)
	t.Ratio = values[7].(float64)
	t.UploadRate = values[8].(float64)
	t.DownloadRate = values[9].(float64)
	t.Finish = values[10].(float64)
	t.Label = values[11].(string)
	t.PeersConnected = values[12].(float64)
	t.PeersTotal = values[13].(float64)
	t.SeedsConnected = values[14].(float64)
	t.SeedsTotal = values[15].(float64)
	t.Availability = values[16].(float64)
	t.Queue = values[17].(float64)
	t.Remaining = values[18].(float64)
	t.State = values[21].(string)
}

func (t *Torrent) SizeMb() float64 {
	return t.Size / 1000000
}
