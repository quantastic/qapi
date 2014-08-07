package qapi

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

// @TODO Implement only accepting /outputting UTC JSON.
type UTCTime struct {
	time.Time
}

func NewTime() *Time {
	return &Time{Id: newId(), Start: time.Now()}
}

type Time struct {
	Id       string     `json:"id"`
	Category *Category  `json:"category"`
	End      *time.Time `json:"end"`
	Start    time.Time  `json:"start"`
	Note     string     `json:"note"`
}

func (t *Time) Duration() time.Duration {
	var end time.Time
	if t.End != nil {
		end = *t.End
	} else {
		end = time.Now()
	}
	return end.Sub(t.Start)
}

type Category struct {
	Name []string `json:"note"`
}

func newId() string {
	b := make([]byte, 16)
	if _, err := io.ReadAtLeast(rand.Reader, b, len(b)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
