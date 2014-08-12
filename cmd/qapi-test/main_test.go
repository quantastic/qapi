package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Globals (local state is preferable, so this list should stay small)
var (
	config Config
)

type Config struct {
	Addr string
	URL  string
}

func ParseFlags(args []string, c *Config) {
	flag.StringVar(&c.Addr, "addr", ":8081", "Http addr to listen on.")
	flag.StringVar(&c.URL, "url", "", "Http url to test against.")
	flag.StringVar(&c.URL, "-v", "", "Http url to test against.")
	flag.Parse()
}

// init() builds and launches a qapi process to run test against. Initially I
// wanted qapi-test to be a regular binary with a main(), but then I realized
// I'd have to come up with my own *testing.T style functionality and decided
// to be lazy for now and just piggy back on *testing.T. Maybe I'll revist this
// in the future.
func init() {
	ParseFlags(os.Args[1:], &config)
	fmt.Printf("-> Running unit tests\n")
	test := exec.Command("go", "test", "github.com/quantastic/qapi")
	testOut, err := test.CombinedOutput()
	if err != nil {
		fatalf("Unit tests failed: %s:\n%s", err, testOut)
	}
	fmt.Printf("-> Building qapi\n")
	tmpDir, err := ioutil.TempDir("", "qapi-test")
	if err != nil {
		fatalf("TempDir: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	binPath := filepath.Join(tmpDir, "qapi")
	build := exec.Command("go", "build", "-o", binPath, "github.com/quantastic/qapi/cmd/qapi")
	if err := build.Run(); err != nil {
		fatalf("%s", err)
	}
	defer os.Remove(binPath)
	// @TODO lsof may fail because its not present, would be nice to detect that.
	lsof := exec.Command("lsof", "-t", "-i", config.Addr)
	pidStr, err := lsof.CombinedOutput()
	if err == nil {
		var pid int
		if _, err := fmt.Sscanf(string(pidStr), "%d", &pid); err != nil {
			fatalf("Bad pid: %s: %s\n", pidStr, err)
		}
		fmt.Printf("-> Killing qapi %d\n", pid)
		proc, err := os.FindProcess(pid)
		if err != nil {
			fatalf("FindProcess: %d: %s", pid, err)
		}
		if err := proc.Kill(); err != nil {
			fatalf("Kill: %d: %s", pid, err)
		}
	}
	fmt.Printf("-> Starting qapi\n")
	config.URL = fmt.Sprintf("http://%s", config.Addr)
	qapi := exec.Command(binPath, "-addr="+config.Addr)
	out := &bytes.Buffer{}
	qapi.Stdout = out
	qapi.Stderr = out
	if err := qapi.Start(); err != nil {
		fatalf("%s", err)
	}
	go func() {
		// @TODO Sometimes qapi doesn't exit (e.g. when qapi-test panics), so we need
		// a way to kill it. Probably using lsof.
		if err := qapi.Wait(); err != nil {
			fatalf("%s\nqapi died: %s", out, err)
		}
	}()
	timeout := time.Now().Add(3 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		var res *http.Response
		res, err = http.Get(config.URL)
		if err != nil {
			continue
		}
		res.Body.Close()
		if res.StatusCode != 200 {
			err = fmt.Errorf("Bad status: %d - %s", res.StatusCode, res.Status)
		}
		break
	}
	if err != nil {
		fatalf("%s", err)
	}
}
