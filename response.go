package flame

type Response struct {
	Build    float64
	Torrents []Torrent
	CacheId  string
}

func (r *Response) Load(data *map[string]interface{}) {
	value := *data
	r.Build = value["build"].(float64)
	r.CacheId = value["torrentc"].(string)
	for _, t := range value["torrents"].([]interface{}) {
		//fmt.Println(t)
		torrent := Torrent{}
		torrent.Load(t.([]interface{}))
		r.Torrents = append(r.Torrents, torrent)
	}
}
