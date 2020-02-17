package utorrent

type Status struct {
	value int
	Started,
	Checking,
	StartAfterCheck,
	Checked,
	Error,
	Paused,
	Queued,
	Loaded bool
}

func NewStatus(val int) *Status {
	return &Status{
		value:           val,
		Started:         (val & (1 << 0)) != 0,
		Checking:        (val & (1 << 1)) != 0,
		StartAfterCheck: (val & (1 << 2)) != 0,
		Checked:         (val & (1 << 3)) != 0,
		Error:           (val & (1 << 4)) != 0,
		Paused:          (val & (1 << 5)) != 0,
		Queued:          (val & (1 << 6)) != 0,
		Loaded:          (val & (1 << 7)) != 0,
	}
}

//func (s *Status) String() string {
//
//}
