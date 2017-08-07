package io_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	//. "io"
	theio "github.com/javouhey/seneca/io"
	"github.com/stretchr/testify/assert"
)

var (
	mapd    map[string]time.Duration
	maps    map[string]theio.VideoSize
	mapf    map[string]float32
	streams []string
)

// fixtures
func init() {
	conv := time.ParseDuration
	d1, _ := conv("00h05m06s")
	d2, _ := conv("00h08m20s")
	mapd = map[string]time.Duration{
		"Duration: 00:05:06.00, start: 0.000000, bitrate: 342 kb/s ": d1,
		"Duration: 00:08:20.53, start: 0.000000, bitrate: 709 kb/s":  d2,
	}

	streams = []string{
		"Stream #0:0: Video: flv1, yuv420p, 426x240, 276 kb/s, 25 tbr, ",
		"Stream #0:0(und): Video: h264 (High) (avc1 / 0x31637661), yuv420p, 960x720 [SAR 1:1 DAR 4:3], 2020 kb/s, 29.97 fps, 29.97 tbr, 30k tbn, ",
		"Stream #0:0(und): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 1280x720, 1037 kb/s, 23.97 fps, 23.97 tbr,",
		"Stream #0:0(und): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 1280x720 , 1037 kb/s, 23.97 fps, 23.97 tbr,",
		"Stream #0:0(und): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 1280x720  , 1037 kb/s, 23.97 fps, 23.97 tbr,",
	}

	maps = map[string]theio.VideoSize{
		streams[0]: theio.VideoSize{426, 240},
		streams[1]: theio.VideoSize{960, 720},
		streams[2]: theio.VideoSize{1280, 720},
		streams[3]: theio.VideoSize{1280, 720},
		streams[4]: theio.VideoSize{1280, 720},
	}

	f1, _ := strconv.ParseFloat("25", 32)
	f2, _ := strconv.ParseFloat("29.97", 32)
	f3, _ := strconv.ParseFloat("23.97", 32)
	mapf = map[string]float32{
		streams[0]: float32(f1),
		streams[1]: float32(f2),
		streams[2]: float32(f3),
	}
}

func TestFps(t *testing.T) {
	assert.False(t, theio.RegexFps1.MatchString(streams[0]))

	// --- fps ---
	assert.True(t, theio.RegexFps1.MatchString(streams[1]))
	matched := theio.RegexFps1.ReplaceAllString(streams[1],
		fmt.Sprintf("${%s}", theio.RegexFps1.SubexpNames()[2]))
	assert.Equal(t, matched, "29.97 fps,", "")

	assert.True(t, theio.RegexFps1.MatchString(streams[2]))
	matched = theio.RegexFps1.ReplaceAllString(streams[2],
		fmt.Sprintf("${%s}", theio.RegexFps1.SubexpNames()[2]))
	assert.Equal(t, matched, "23.97 fps,", "")

	// --- tbr ---
	assert.True(t, theio.RegexFps2.MatchString(streams[0]))
	matched = theio.RegexFps2.ReplaceAllString(streams[0],
		fmt.Sprintf("${%s}", theio.RegexFps2.SubexpNames()[2]))
	assert.Equal(t, matched, "25 tbr,", "")

	assert.True(t, theio.RegexFps2.MatchString(streams[2]))
	matched = theio.RegexFps2.ReplaceAllString(streams[2],
		fmt.Sprintf("${%s}", theio.RegexFps2.SubexpNames()[2]))
	assert.Equal(t, matched, "23.97 tbr,", "")

	for k, v := range mapf {
		if d, _ := theio.ParseFps(k); d != v {
			t.Errorf("tsk tsk: Fps")
		}
	}
}

func TestSizeRegex(t *testing.T) {
	for k, v := range maps {
		d, _ := theio.ParseDimension(k)
		if d != v {
			t.Errorf("tsk tsk: Size")
		}
	}
}

func TestDurationRegex(t *testing.T) {
	_, err := theio.ParseDuration("")
	assert.Error(t, err, "")
	assert.Equal(t, err, theio.InvalidDuration, "empty duration should fail")

	for k, v := range mapd {
		if d, _ := theio.ParseDuration(k); d != v {
			t.Errorf("tsk tsk: Duration")
		}
	}
}
