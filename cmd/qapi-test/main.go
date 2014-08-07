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

func main() {
	message, err := ParseFlags(os.Args[1:], &config)
	if err == flag.ErrHelp {
		fmt.Printf("%s\n", message)
		os.Exit(0)
	} else if err != nil {
		fatalf("%s", message)
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
	fmt.Printf("-> Starting qapi\n")
	config.URL = fmt.Sprintf("http://%s", config.Addr)
	qapi := exec.Command(binPath, "-addr="+config.Addr)
	out := &bytes.Buffer{}
	qapi.Stdout = out
	qapi.Stderr = out
	if err := qapi.Start(); err != nil {
		fatalf("%s", err)
	}
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
	fmt.Printf("-> Running http tests\n")
	if err := qapi.Wait(); err != nil {
		fatalf("%s\nqapi died: %s", out, err)
	}
}
