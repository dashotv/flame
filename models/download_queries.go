package models

func (s *DownloadStore) Active() ([]Download, error) {
	q := s.Query()
	return q.In("status", []string{"searching", "loading", "managing", "downloading", "reviewing"}).Run()
}
