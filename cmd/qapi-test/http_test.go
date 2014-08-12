package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

func Test_Home(t *testing.T) {
	res, err := http.Get(config.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: %d", res.StatusCode)
	}
}

func Test_CreateTime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	minLater := now.Add(time.Minute)
	tests := []struct {
		Data   interface{}
		Status int
		Result interface{}
	}{
		{
			Data:   &Time{},
			Status: http.StatusBadRequest,
			Result: map[string]string{
				"error":   "bad request",
				"message": "Required field is missing: category",
			},
		},
		{
			Data:   &Time{Category: []string{"A", "B"}},
			Status: http.StatusBadRequest,
			Result: map[string]string{
				"error":   "bad request",
				"message": "Required field is missing: start",
			},
		},
		{
			Data:   &Time{Category: []string{"A", "B"}, Start: &now},
			Status: http.StatusCreated,
			Result: func(result map[string]interface{}) error {
				r := result["time"].(map[string]interface{})
				start, err := time.Parse(time.RFC3339, r["start"].(string))
				if err != nil {
					return err
				} else if !start.Equal(now) {
					return fmt.Errorf("Want start=%s, got: %s", now, start)
				} else if start.Location() != time.UTC {
					return fmt.Errorf("Want start location UTC, got: %s", start.Location())
				}
				want := []string{"A", "B"}
				category := r["category"].(map[string]interface{})
				if diff := pretty.Compare(category["name"], want); diff != "" {
					return fmt.Errorf("Bad category:\n%s", diff)
				}
				if r["end"] != nil {
					return fmt.Errorf("Want end=nil, got: %s", r["end"])
				}
				if r["note"].(string) != "" {
					return fmt.Errorf("Want empty note, got: %s", r["note"])
				}
				return nil
			},
		},
		{
			Data: &Time{
				Category: []string{"A", "B"},
				Start:    &now,
				End:      &minLater,
				Note:     "Hello World",
			},
			Status: http.StatusCreated,
			Result: func(result map[string]interface{}) error {
				r := result["time"].(map[string]interface{})
				end, err := time.Parse(time.RFC3339, r["end"].(string))
				if err != nil {
					return err
				} else if !end.Equal(minLater) {
					return fmt.Errorf("Want end=%s, got: %s", now, end)
				} else if end.Location() != time.UTC {
					return fmt.Errorf("Want end location UTC, got: %s", end.Location())
				}
				if r["note"].(string) != "Hello World" {
					return fmt.Errorf("Want note=Hello World, got: %s", r["note"])
				}
				return nil
			},
		},
	}
	for i, test := range tests {
		var result map[string]interface{}
		res, body, err := jsonPost("/times", test.Data, &result)
		if err != nil {
			t.Error(err)
			continue
		}
		if res.StatusCode != test.Status {
			t.Errorf("test %d: Bad status %d: %s", i, res.StatusCode, body)
			continue
		}
		if resultFunc, ok := test.Result.(func(map[string]interface{}) error); ok {
			if err := resultFunc(result); err != nil {
				t.Errorf("test %d: %s", i, err)
				continue
			}
		} else if diff := pretty.Compare(result, test.Result); diff != "" {
			t.Errorf("test %d:\n%s", i, diff)
			continue
		}
	}
}

type Time struct {
	Id       string     `json:"id,omitempty"`
	Category []string   `json:"category,omitempty"`
	End      *time.Time `json:"end,omitempty"`
	Start    *time.Time `json:"start,omitempty"`
	Note     string     `json:"note,omitempty"`
}

type UTCTime struct {
	time.Time
}

func (u UTCTime) MarshalJSON() ([]byte, error) {
	return u.UTC().Truncate(time.Second).MarshalJSON()
}
