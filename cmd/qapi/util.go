package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func jsonDecode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	const maxSize = 1024 * 1024
	d := json.NewDecoder(io.LimitReader(r.Body, maxSize))
	if err := d.Decode(v); err != nil {
		jsonWrite(w, http.StatusBadRequest, map[string]string{
			"error":   "bad request",
			"message": err.Error(),
		})
		return err
	}
	return nil
}

func jsonWrite(w http.ResponseWriter, status int, v interface{}) error {
	if err, ok := v.(error); ok {
		v = map[string]string{"error": err.Error()}
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - %s", err)
		return err
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
