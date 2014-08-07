package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Addr string
	URL  string
}

func ParseFlags(args []string, c *Config) (string, error) {
	msg := bytes.NewBuffer(nil)
	f := flag.NewFlagSet("qapi-test", flag.ContinueOnError)
	f.SetOutput(msg)
	f.Usage = func() {}
	f.StringVar(&c.Addr, "addr", ":8081", "Http addr to listen on.")
	f.StringVar(&c.URL, "url", "", "Http url to test against.")
	fmt.Fprintf(msg, "qapi-test")
	f.PrintDefaults()
	fmt.Fprintf(msg, "\n")
	err := f.Parse(args)
	if err != nil {
		return strings.TrimSpace(msg.String()), err
	}
	return "", nil
}
