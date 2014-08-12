package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/quantastic/qapi"
)

// Globals (local state is preferable, so this list should stay small)
var (
	config Config
	db     *qapi.Db
)

func main() {
	message, err := ParseFlags(os.Args[1:], &config)
	if err == flag.ErrHelp {
		fmt.Printf("%s\n", message)
		os.Exit(0)
	} else if err != nil {
		fatalf("%s", message)
	}
	db, err = qapi.OpenDb(config.Dir)
	if err != nil {
		fatalf("OpenDb: %s", err)
	}
	server := &http.Server{Addr: config.Addr, Handler: Router()}
	if err := server.ListenAndServe(); err != nil {
		fatalf("ListenAndServe: %s", err)
	}
}
