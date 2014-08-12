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
	db := &Db{dir: dir, times: make(map[string]Time)}
	if err := db.loadTimes(); err != nil {
		return nil, err
	}
	return db, nil
}

type Db struct {
	dir   string
	times map[string]Time
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
		times = append(times, t)
	}
	slice.Sort(times, func(i, j int) bool {
		return times[i].Start.After(times[j].Start)
	})
	return times, nil
}

func (d *Db) SaveTime(t *Time) error {
	now := time.Now().UTC().Truncate(time.Second)
	if t.Id == "" {
		t.Id = mustUUID()
		t.Created = now
	}
	t.Updated = now
	t.normalize()
	d.times[t.Id] = *t
	return d.saveTimes()
}

func mustUUID() string {
	b := make([]byte, 16)
	if _, err := io.ReadAtLeast(rand.Reader, b, len(b)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
