package io

import (
	"testing"
    "time"
    //"log"
    //"strings"

    theio "github.com/javouhey/seneca/io"
)

// fixtures
var (
    mapd map[string]time.Duration
)

func init() {
    conv := time.ParseDuration
    d1, _ := conv("00h05m06s")
    d2, _ := conv("00h08m20s")
    mapd = map[string]time.Duration{
      "Duration: 00:05:06.00, start: 0.000000, bitrate: 342 kb/s ": d1,
      "Duration: 00:08:20.53, start: 0.000000, bitrate: 709 kb/s": d2,
    }
}

func TestDurationRegex(t *testing.T) {
    _, err := theio.ParseDuration("")
    if err != theio.InvalidDuration {
        t.Errorf("empty duration should fail")
    }

    for k, v := range mapd {
        if d, _ := theio.ParseDuration(k); d != v {
            t.Errorf("tsk tsk")
        }
    }
}
