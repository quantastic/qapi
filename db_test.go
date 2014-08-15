package qapi

import (
	"fmt"
	"os"

	"github.com/kylelemons/godebug/pretty"
)
import "io/ioutil"
import (
	"testing"
	"time"
)

func Test_Db_SaveTime(t *testing.T) {
	offset, err := time.Parse(time.RFC3339, "2014-08-13T09:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		Input []Time
		Want  []Time
	}{
		// new active entry sets end of old active entry
		{
			Input: []Time{
				{
					Category: []string{"A"},
					Start:    offset.Add(-10 * time.Minute),
					End:      nil,
				},
				{
					Category: []string{"B"},
					Start:    offset.Add(-5 * time.Minute),
					End:      nil,
				},
			},
			Want: []Time{
				{
					Category: []string{"B"},
					Start:    offset.Add(-5 * time.Minute),
					End:      nil,
				},
				{
					Category: []string{"A"},
					Start:    offset.Add(-10 * time.Minute),
					End:      timeAddr(offset.Add(-5 * time.Minute)),
				},
			},
		},
		// if new active start is before old active start, old active end is set to
		// old active start (duration=0).
		{
			Input: []Time{
				{
					Category: []string{"A"},
					Start:    offset.Add(-10 * time.Minute),
					End:      nil,
				},
				{
					Category: []string{"B"},
					Start:    offset.Add(-15 * time.Minute),
					End:      nil,
				},
			},
			Want: []Time{
				{
					Category: []string{"B"},
					Start:    offset.Add(-15 * time.Minute),
					End:      nil,
				},
				{
					Category: []string{"A"},
					Start:    offset.Add(-10 * time.Minute),
					End:      timeAddr(offset.Add(-10 * time.Minute)),
				},
			},
		},
	}
	for i, test := range tests {
		dir, err := ioutil.TempDir("", "db_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		db, err := OpenDb(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, input := range test.Input {
			if err := db.SaveTime(&input); err != nil {
				t.Fatal(err)
			}
		}
		results, err := db.Times()
		if err != nil {
			t.Fatal(err)
		}
		for j, result := range results {
			want := test.Want[j]
			if diff := pretty.Compare(result.Category, want.Category); diff != "" {
				t.Errorf("test %d.%d: bad category:\n%s", i, j, diff)
			}
			if diff := pretty.Compare(result.Start.String(), want.Start.String()); diff != "" {
				t.Errorf("test %d.%d: bad start:\n%s", i, j, diff)
			}
			if diff := pretty.Compare(fmt.Sprintf("%s", result.End), fmt.Sprintf("%s", want.End)); diff != "" {
				t.Errorf("test %d.%d: bad end:\n%s", i, j, diff)
			}
		}
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}
}

func TestShadowTimes(t *testing.T) {
	offset, err := time.Parse(time.RFC3339, "2014-08-13T09:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		Input []Time
		Want  []Time
	}{
		// B shadows parts of A
		{
			Input: []Time{
				{
					Category: []string{"B"},
					End:      timeAddr(offset.Add(-5 * time.Minute)),
					Start:    offset.Add(-15 * time.Minute),
				},
				{
					Category: []string{"A"},
					End:      timeAddr(offset.Add(-10 * time.Minute)),
					Start:    offset.Add(-20 * time.Minute),
				},
			},
			Want: []Time{
				{
					Category: []string{"B"},
					End:      timeAddr(offset.Add(-5 * time.Minute)),
					Start:    offset.Add(-15 * time.Minute),
				},
				{
					Category: []string{"A"},
					End:      timeAddr(offset.Add(-15 * time.Minute)),
					Start:    offset.Add(-20 * time.Minute),
				},
			},
		},
		// B shadows all of A
		{
			Input: []Time{
				{
					Category: []string{"B"},
					End:      timeAddr(offset.Add(-5 * time.Minute)),
					Start:    offset.Add(-20 * time.Minute),
				},
				{
					Category: []string{"A"},
					End:      timeAddr(offset.Add(-10 * time.Minute)),
					Start:    offset.Add(-15 * time.Minute),
				},
			},
			Want: []Time{
				{
					Category: []string{"B"},
					End:      timeAddr(offset.Add(-5 * time.Minute)),
					Start:    offset.Add(-20 * time.Minute),
				},
				{
					Category: []string{"A"},
					End:      timeAddr(offset.Add(-15 * time.Minute)),
					Start:    offset.Add(-15 * time.Minute),
				},
			},
		},
	}
	for i, test := range tests {
		results := ShadowTimes(test.Input)
		if len(results) != len(test.Want) {
			t.Errorf("test %d: want %d results, got %d", i, len(test.Want), len(results))
			continue
		}
		for j, result := range results {
			want := test.Want[j]
			if diff := pretty.Compare(result.Category, want.Category); diff != "" {
				t.Errorf("test %d.%d: bad category:\n%s", i, j, diff)
			}
			if diff := pretty.Compare(result.Start.String(), want.Start.String()); diff != "" {
				t.Errorf("test %d.%d: bad start:\n%s", i, j, diff)
			}
			if diff := pretty.Compare(fmt.Sprintf("%s", result.End), fmt.Sprintf("%s", want.End)); diff != "" {
				t.Errorf("test %d.%d: bad end:\n%s", i, j, diff)
			}
		}
	}
}

func timeAddr(t time.Time) *time.Time {
	return &t
}
