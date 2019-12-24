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
	t.Hash = t.getString(values[0])
	t.Status = int(t.getFloat64(values[1]))
	t.Name = t.getString(values[2])
	t.Size = t.getFloat64(values[3])
	t.Progress = t.getFloat64(values[4]) / 10
	t.Downloaded = t.getFloat64(values[5])
	t.Uploaded = t.getFloat64(values[6])
	t.Ratio = t.getFloat64(values[7])
	t.UploadRate = t.getFloat64(values[8])
	t.DownloadRate = t.getFloat64(values[9])
	t.Finish = t.getFloat64(values[10])
	t.Label = t.getString(values[11])
	t.PeersConnected = t.getFloat64(values[12])
	t.PeersTotal = t.getFloat64(values[13])
	t.SeedsConnected = t.getFloat64(values[14])
	t.SeedsTotal = t.getFloat64(values[15])
	t.Availability = t.getFloat64(values[16])
	t.Queue = t.getFloat64(values[17])
	t.Remaining = t.getFloat64(values[18])
	t.State = t.getString(values[21])
}

func (t *Torrent) getString(value interface{}) string {
	if value != nil {
		return value.(string)
	}
	return ""
}

func (t *Torrent) getFloat64(value interface{}) float64 {
	if value != nil {
		return value.(float64)
	}
	return 0.0
}

func (t *Torrent) SizeMb() float64 {
	return t.Size / 1000000
}
