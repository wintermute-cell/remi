package model

import "time"

type Reminder struct {
	Timestamp Timestamp
	Message   string
}

type Timestamp struct {
	DateTime time.Time
	HasTime  bool      // Is the timestamp time relevant or just the date?
	ShowFrom time.Time // When should the reminder start showing?
}

func (t Timestamp) String() string {
	if t.HasTime {
		return t.DateTime.Format("02.01.06@15:04")
	}
	return t.DateTime.Format("02.01.06      ")
}
