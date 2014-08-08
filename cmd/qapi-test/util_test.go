package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func jsonPost(path, data, result interface{}) (*http.Response, string, error) {
	reqBody := &bytes.Buffer{}
	e := json.NewEncoder(reqBody)
	if err := e.Encode(data); err != nil {
		return nil, "", err
	}
	res, err := http.Post(config.URL+"/times", "application/json", reqBody)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	resBody := &bytes.Buffer{}
	if _, err := io.Copy(resBody, res.Body); err != nil {
		return nil, "", err
	}
	if err := json.Unmarshal(resBody.Bytes(), result); err != nil {
		return nil, "", fmt.Errorf("%s: %s", err, resBody)
	}
	return res, resBody.String(), nil
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
