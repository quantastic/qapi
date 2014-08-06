package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/quantastic/api"
)

var (
	addr = flag.String("addr", ":8080", "Http host:port to listen on.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "qapid [flags]\n\nAvailable flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	handler := httpHandler()
	server := &http.Server{Addr: *addr, Handler: handler}
	if err := server.ListenAndServe(); err != nil {
		fatalf("ListenAndServe: %s", err)
	}
}

func httpHandler() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/times", func(w http.ResponseWriter, r *http.Request) {
		t := &api.Time{}
		if err := jsonDecode(w, r, t); err != nil {
			return
		}
		jsonWrite(w, http.StatusOK, map[string]interface{}{"time": t})
	}).Methods("POST")
	return mux
}

func jsonWrite(w http.ResponseWriter, status int, v interface{}) error {
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

func jsonDecode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	const maxSize = 1024 * 1024
	d := json.NewDecoder(io.LimitReader(r.Body, maxSize))
	if err := d.Decode(v); err != nil {
		jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "bad json"})
		return err
	}
	return nil
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
