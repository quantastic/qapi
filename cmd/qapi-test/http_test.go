package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
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
	reqBody := &bytes.Buffer{}
	e := json.NewEncoder(reqBody)
	if err := e.Encode(map[string]interface{}{}); err != nil {
		t.Fatal(err)
	}
	res, err := http.Post(config.URL+"/times", "application/json", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	resBody := &bytes.Buffer{}
	if _, err := io.Copy(resBody, res.Body); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Bad status %d: %s", res.StatusCode, resBody)
	}
}
