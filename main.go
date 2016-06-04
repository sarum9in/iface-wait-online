package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var ip = flag.String("ip", "ip", "iproute2 path")
var size = flag.Int("size", 1, "Expected size of response")
var timeout = flag.Int("timeout", 60, "Number of seconds to try before failure")

var ErrRetry = errors.New("Retry")

const (
	retryInterval = time.Millisecond * 100
)

func doRequest() (string, error) {
	cmd := exec.Command(*ip, append([]string{"-oneline"}, flag.Args()...)...)
	output, err := cmd.Output()
	if err != nil {
		if werr, ok := err.(*exec.ExitError); ok {
			log.Printf("Failed: %s, retrying...",
				strings.TrimRight(string(werr.Stderr), "\n"))
			return "", ErrRetry
		} else {
			return "", err
		}
	}
	return string(output), nil
}

func waitForRequest() error {
	tryUntil := time.Now().Add(time.Duration(*timeout) * time.Second)
	for time.Now().Before(tryUntil) {
		output, err := doRequest()
		if err == nil {
			c := strings.Count(output, "\n")
			if c != *size {
				fmt.Print(".")
			} else {
				fmt.Println("+")
				return nil
			}
		} else if err != ErrRetry {
			return err
		}
		time.Sleep(retryInterval)
	}
	fmt.Println()
	return fmt.Errorf("Timeout")
}

func main() {
	flag.Parse()
	if err := waitForRequest(); err != nil {
		log.Fatal(err)
	}
}
