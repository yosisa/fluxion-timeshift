package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	since  = kingpin.Flag("since", "").Short('s').String()
	until  = kingpin.Flag("until", "").Short('u').String()
	dur    = kingpin.Flag("duration", "").Short('d').String()
	tag    = kingpin.Flag("tag", "").Short('t').String()
	server = kingpin.Arg("url", "").Required().String()
)

func parseDuration(since, until, duration string) (s int64, u int64, err error) {
	if since == "" && until == "" {
		return
	}
	if since != "" && until != "" {
		if s, err = parseTime(since); err != nil {
			return
		}
		if u, err = parseTime(until); err != nil {
			return
		}
	}

	var delta time.Duration
	if duration != "" {
		if delta, err = time.ParseDuration(duration); err != nil {
			return
		}
	}
	if since != "" {
		if s, err = parseTime(since); err != nil {
			return
		}
		if delta > 0 {
			u = s + int64(delta.Seconds())
		}
		return
	}

	if u, err = parseTime(until); err != nil {
		return
	}
	if delta > 0 {
		s = u - int64(delta.Seconds())
	}
	return
}

func parseTime(t string) (int64, error) {
	d, err := time.ParseDuration(t)
	if err == nil {
		return time.Now().Add(d).Unix(), nil
	}
	return 0, fmt.Errorf("Unsupported time format: %s", t)
}

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()

	s, u, err := parseDuration(*since, *until, *dur)
	if err != nil {
		log.Fatal(err)
	}
	var params []string
	if s != 0 {
		params = append(params, fmt.Sprintf("since=%d", s))
	}
	if u != 0 {
		params = append(params, fmt.Sprintf("until=%d", u))
	}
	if *tag != "" {
		params = append(params, "tag="+*tag)
	}

	resp, err := http.Get(*server + "?" + strings.Join(params, "&"))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
