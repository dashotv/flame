package config

type Config struct {
	Utorrent struct {
		URL string
	}
	Nzbget struct {
		URL string
	}
	Port int
	Mode string
}
