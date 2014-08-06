package api

import "time"

// @TODO Implement only accepting /outputting UTC JSON.
type UTCTime struct {
	time.Time
}

type Time struct {
	Id    string   `json:"id"`
	End   *UTCTime `json:"end"`
	Start UTCTime  `json:"start"`
	Note  string   `json:"note"`
}
