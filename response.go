package flame

type Response struct {
	Build    float64
	Torrents []Torrent
	CacheId  string
}

func (r *Response) Load(data *map[string]interface{}) {
	value := *data
	r.Build = value["build"].(float64)

	if val, ok := value["torrentc"]; ok {
		r.CacheId = val.(string)
	}

	if val, ok := value["torrents"]; ok {
		for _, t := range val.([]interface{}) {
			//fmt.Println(t)
			torrent := Torrent{}
			torrent.Load(t.([]interface{}))
			r.Torrents = append(r.Torrents, torrent)
		}
	}
}

func (r *Response) Count() (int) {
	return len(r.Torrents)
}
