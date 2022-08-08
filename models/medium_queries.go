package models

import "time"

func (s *MediumStore) Upcoming() ([]Medium, error) {
	q := s.Query()
	return q.Where("missing", nil).LessThanEqual("release_date", time.Now()).Run()
}
