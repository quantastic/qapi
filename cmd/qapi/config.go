package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Addr string
}

func ParseFlags(args []string, c *Config) (string, error) {
	msg := bytes.NewBuffer(nil)
	f := flag.NewFlagSet("qapi", flag.ContinueOnError)
	f.SetOutput(msg)
	f.Usage = func() {}
	f.StringVar(&c.Addr, "addr", ":8080", "Http host:port to listen on.")
	fmt.Fprintf(msg, "qapi")
	f.PrintDefaults()
	fmt.Fprintf(msg, "\n")
	err := f.Parse(args)
	if err != nil {
		return strings.TrimSpace(msg.String()), err
	}
	return "", nil
}
