package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Addr string
	Dir  string
	URL  string
}

func ParseFlags(args []string, c *Config) (string, error) {
	msg := bytes.NewBuffer(nil)
	f := flag.NewFlagSet("qapi", flag.ContinueOnError)
	f.SetOutput(msg)
	f.Usage = func() {}
	f.StringVar(&c.Addr, "addr", ":8080", "Http host:port to listen on.")
	f.StringVar(&c.Dir, "dir", "qdata", "Path to dir to store data in.")
	// @TODO Adjust url according to addr if default is used
	f.StringVar(&c.URL, "url", "http://localhost:8080", "Base url to use for links.")
	fmt.Fprintf(msg, "qapi")
	f.PrintDefaults()
	fmt.Fprintf(msg, "\n")
	err := f.Parse(args)
	if err != nil {
		return strings.TrimSpace(msg.String()), err
	}
	return "", nil
}
