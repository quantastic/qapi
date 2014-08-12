package qapi

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"
)

func NewTime() Time {
	return Time{Start: time.Now()}
}

type Time struct {
	Id       string
	Category []string
	End      *time.Time
	Start    time.Time
	Note     string
	Created  time.Time
	Updated  time.Time
}

func (t Time) Duration() time.Duration {
	var end time.Time
	if t.End != nil {
		end = *t.End
	} else {
		end = time.Now()
	}
	return end.Sub(t.Start)
}

func (t Time) Valid() error {
	if t.Start.IsZero() {
		return fmt.Errorf("Start must not be zero")
	}
	if t.End != nil && t.End.Before(t.Start) {
		return fmt.Errorf("End must not be before Start.")
	}
	if len(t.Category) < 1 {
		return fmt.Errorf("Category must not be empty.")
	}
	for _, part := range t.Category {
		if strings.TrimSpace(part) == "" {
			return fmt.Errorf("Category part must not be empty.")
		}
	}
	return nil
}

func (t *Time) normalize() {
	t.Start = t.Start.UTC().Truncate(time.Second)
	if t.End != nil {
		ue := t.End.UTC().Truncate(time.Second)
		t.End = &ue
	}
}

func newId() string {
	b := make([]byte, 16)
	if _, err := io.ReadAtLeast(rand.Reader, b, len(b)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
