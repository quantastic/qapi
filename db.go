package qapi

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bradfitz/slice"
)

const timeFile = "time.json"

func OpenDb(dir string) (*Db, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}
	db := &Db{dir: dir, times: make(map[string]*Time)}
	if err := db.loadTimes(); err != nil {
		return nil, err
	}
	return db, nil
}

// Db is a simple/naive file based data storage. It will eventually be replaced
// with adapters for real databases, but for now this is convenient to
// prototype with.
type Db struct {
	dir   string
	times map[string]*Time
}

func (d *Db) loadTimes() error {
	data, err := ioutil.ReadFile(filepath.Join(d.dir, timeFile))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &d.times); err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) saveTimes() error {
	file, err := ioutil.TempFile("", timeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	if err := e.Encode(d.times); err != nil {
		return err
	}
	return os.Rename(file.Name(), filepath.Join(d.dir, timeFile))
}

func (d *Db) Times() ([]Time, error) {
	times := make([]Time, 0, len(d.times))
	for _, t := range d.times {
		t.normalize()
		times = append(times, t.Copy())
	}
	Sort(times)
	return times, nil
}

func (d *Db) active() *Time {
	var active *Time
	for _, t := range d.times {
		if t.End == nil {
			// @TODO Return the latest active entry and make sure other methods
			// only return a single active entry if this kind of data corruption
			// is present. Then also provide a method to list the corrupted entries.
			if active != nil {
				panic("More than one active entry")
			}
			active = t
		}
	}
	return active
}

func (d *Db) SaveTime(t *Time) error {
	c := t.Copy()
	t = &c
	now := time.Now().UTC().Truncate(time.Second)
	if t.Id == "" {
		t.Id = mustUUID()
		t.Created = now
	}
	t.Updated = now
	t.normalize()
	if t.End == nil {
		active := d.active()
		if active != nil {
			var s time.Time
			if t.Start.Before(active.Start) {
				s = active.Start
			} else {
				s = t.Start
			}
			active.End = &s
			active.Updated = now
		}
	}
	d.times[t.Id] = t
	return d.saveTimes()
}

// @TODO Write test
func Sort(times []Time) {
	slice.Sort(times, func(i, j int) bool {
		if times[i].End == nil {
			// i is active entry, sort it before j
			return true
		} else if times[j].End == nil {
			// j is active entry, sort it before i
			return false
		}
		if times[i].End.Equal(*times[j].End) {
			// i and j have same end, sort i before j if i started after j
			return times[i].Start.After(times[j].Start)
		}
		// sort i before j if it ends after j
		return times[i].End.After(*times[j].End)
	})
}

func Shadow(times []Time) []Time {
	results := make([]Time, 0, len(times))
	var prev *Time
	for _, t := range times {
		c := t.Copy()
		if prev != nil {
			if c.End.After(prev.Start) {
				e := prev.Start
				if e.Before(c.Start) {
					e = c.Start
				}
				c.End = &e
			}
		}
		results = append(results, c)
		prev = &c
	}
	return results
}

func mustUUID() string {
	b := make([]byte, 16)
	if _, err := io.ReadAtLeast(rand.Reader, b, len(b)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
