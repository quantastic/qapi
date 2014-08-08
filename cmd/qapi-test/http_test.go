package main

import (
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
				"message": "Missing start value.",
			},
		},
	}
	for _, test := range tests {
		var result interface{}
		res, body, err := jsonPost("/times", map[string]interface{}{}, &result)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.Status {
			t.Fatalf("Bad status %d: %s", res.StatusCode, body)
		}
		if diff := pretty.Compare(result, test.Result); diff != "" {
			t.Fatalf("%s", diff)
		}
	}
}

type Time struct {
	Id    string
	End   *time.Time
	Start time.Time
	Note  string
}
