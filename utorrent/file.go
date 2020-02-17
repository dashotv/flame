package utorrent

/*
   #, "FILE_NAME": 0
    #, "FILE_SIZE": 1
    #, "FILE_DOWNLOADED": 2
    #, "FILE_PRIORITY": 3
    #, "FILE_FIRST_PIECE": 4
    #, "FILE_NUM_PIECES": 5
    #, "FILE_STREAMABLE": 6
    #, "FILE_ENCODED_RATE": 7
    #, "FILE_DURATION": 8
    #, "FILE_WIDTH": 9
    #, "FILE_HEIGHT": 10
    #, "FILE_STREAM_ETA": 11
    #, "FILE_STREAMABILITY": 12
*/
type File struct {
	Name       string
	Size       float64
	Downloaded float64
	Priority   int
}

func (f *File) Load(values []interface{}) {
	f.Name = getString(values[0])
	f.Size = getFloat64(values[1])
	f.Downloaded = getFloat64(values[2])
	f.Priority = int(getFloat64(values[3]))
}

func (f *File) SizeMb() float64 {
	return f.Size / 1000000
}

func (f *File) DownloadedPercent() float64 {
	return (f.Downloaded / f.Size) * 100
}
