package utorrent

//http://help.utorrent.com/customer/en/portal/articles/1573947-torrent-labels-list---webapi

type Torrent struct {
	Hash           string
	Status         int
	State          string
	Name           string
	Size           float64 // in bytes
	Progress       float64 // float64 in per mils
	Downloaded     float64 // in bytes
	Uploaded       float64 // in bytes
	Ratio          float64 // float64 in per mils
	UploadRate     float64 // float64 in bytes / second
	DownloadRate   float64 // float64 in bytes / second
	Finish         float64 // float64 seconds
	Label          string
	PeersConnected float64
	PeersTotal     float64
	SeedsConnected float64
	SeedsTotal     float64
	Availability   float64 // in 1/65535ths
	Queue          float64
	Remaining      float64 // in bytes
	Path           string
	Files          []*File
}

func (t *Torrent) Load(values []interface{}) {
	t.Hash = getString(values[0])
	t.Status = int(getFloat64(values[1]))
	t.Name = getString(values[2])
	t.Size = getFloat64(values[3])
	t.Progress = getFloat64(values[4]) / 10
	t.Downloaded = getFloat64(values[5])
	t.Uploaded = getFloat64(values[6])
	t.Ratio = getFloat64(values[7])
	t.UploadRate = getFloat64(values[8])
	t.DownloadRate = getFloat64(values[9])
	t.Finish = getFloat64(values[10])
	t.Label = getString(values[11])
	t.PeersConnected = getFloat64(values[12])
	t.PeersTotal = getFloat64(values[13])
	t.SeedsConnected = getFloat64(values[14])
	t.SeedsTotal = getFloat64(values[15])
	t.Availability = getFloat64(values[16])
	t.Queue = getFloat64(values[17])
	t.Remaining = getFloat64(values[18])
	t.State = getString(values[21])
	t.Path = getString(values[26])
}

func (t *Torrent) AddFile(file *File) {
	t.Files = append(t.Files, file)
}

func (t *Torrent) SizeMb() float64 {
	return t.Size / 1000000
}
