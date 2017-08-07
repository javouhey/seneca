/*
Copyright 2014 Gavin Bong.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific
language governing permissions and limitations under the
License.
*/

package io

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/javouhey/seneca/util"
	"labix.org/v2/pipe"
)

const (
	sVideo    = "Video:"
	sDuration = "Duration:"
)

var (
	// Duration: 00:08:20.53, start: 0.000000, bitrate: 709 kb/s
	Regex1 = regexp.MustCompile(`^Duration: (?P<duration>\d{2}:\d{2}:\d{2}).\d{2}(.*)$`)

	// .. yuv420p, 960x720 [SAR 1:1 DAR 4:3], ..
	Regex2 = regexp.MustCompile(`^(?P<prefix>.*?)(?P<size>\d{3,}?x\d{2,}?)([,\s])(?P<postfix>.*)$`)

	// .. , 29.97 fps, 29.97 tbr,
	RegexFps1 = regexp.MustCompile(`^(?P<prefix>.*?)(?P<fps>\d{1,}\.?\d* fps,)(?P<postfix>.*)$`)
	RegexFps2 = regexp.MustCompile(`^(?P<prefix>.*?)(?P<tbr>\d{1,}\.?\d* tbr,)(?P<postfix>.*)$`)

	InvalidDuration  = errors.New("Duration input is invalid")
	InvalidVideoSize = errors.New("Cannot parse for WxH")
	InvalidFps       = errors.New("Cannot parse for fps/tbr")
)

func ParseFps(raw string) (float32, error) {
	empty := float32(0.00)
	if util.IsEmpty(raw) {
		return empty, InvalidFps
	}

	raw = strings.TrimSpace(raw)

	fps1 := RegexFps1.MatchString(raw)
	fps2 := RegexFps2.MatchString(raw)

	if !fps1 && !fps2 {
		return empty, InvalidFps
	}

	if fps1 {
		matched := RegexFps1.ReplaceAllString(raw,
			fmt.Sprintf("${%s}", RegexFps1.SubexpNames()[2]))
		parts := strings.Split(matched, " ")
		f, _ := strconv.ParseFloat(parts[0], 32)
		return float32(f), nil
	}

	if fps2 {
		matched := RegexFps2.ReplaceAllString(raw,
			fmt.Sprintf("${%s}", RegexFps2.SubexpNames()[2]))
		parts := strings.Split(matched, " ")
		f, _ := strconv.ParseFloat(parts[0], 32)
		return float32(f), nil
	}

	return empty, InvalidFps
}

func ParseDimension(raw string) (VideoSize, error) {
	EmptyVid := VideoSize{}
	if util.IsEmpty(raw) {
		return EmptyVid, InvalidVideoSize
	}

	raw = strings.TrimSpace(raw)
	if !Regex2.MatchString(raw) {
		return EmptyVid, InvalidVideoSize
	}

	matched := Regex2.ReplaceAllString(raw,
		fmt.Sprintf("${%s}", Regex2.SubexpNames()[2]))
	parts := strings.Split(matched, "x")
	if len(parts) != 2 {
		return EmptyVid, InvalidVideoSize
	}

	width, err0 := strconv.Atoi(parts[0])
	height, err1 := strconv.Atoi(parts[1])
	retval := VideoSize{uint16(width), uint16(height)}
	if err0 != nil || err1 != nil {
		return EmptyVid, err0
	}
	return retval, nil
}

func ParseDuration(raw string) (time.Duration, error) {
	if util.IsEmpty(raw) {
		return 0, InvalidDuration
	}

	raw = strings.TrimSpace(raw)

	if !Regex1.MatchString(raw) {
		return 0, InvalidDuration
	}

	matched := Regex1.ReplaceAllString(raw,
		fmt.Sprintf("${%s}", Regex1.SubexpNames()[1]))
	parts := strings.Split(matched, ":")
	if len(parts) != 3 {
		return 0, InvalidDuration
	}

	goduration := fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2])
	retval, err := time.ParseDuration(goduration)
	if err != nil {
		return 0, err
	}
	return retval, nil
}

// Parses output from ffprobe
func parse(data *bytes.Buffer) (*VideoReader, error) {
	vid := new(VideoReader)

	p := pipe.Line(
		pipe.Read(bytes.NewReader(data.Bytes())),

		pipe.Filter(func(line []byte) bool {
			s := string(line)
			s = strings.TrimSpace(s)
			return strings.HasPrefix(s, sDuration) ||
				strings.Index(s, sVideo) >= 0
		}),

		Processor(func(line []byte) []byte {
			s := chomp(line)
			if strings.HasPrefix(s, sDuration) {
				d, _ := ParseDuration(s)
				// TODO log the err ??
				vid.Duration = d
				return make([]byte, 0)
			} else {
				return line // leave untouched
			}
		}),

		Processor(func(line []byte) []byte {
			s := chomp(line)
			if strings.Index(s, sVideo) >= 0 {
				dims, _ := ParseDimension(s)
				fps, _ := ParseFps(s)
				// TODO log the err ??
				vid.VideoSize = dims
				vid.Fps = fps
				return make([]byte, 0)
			} else {
				return line // leave untouched
			}
		}),

		CustomWrite(os.Stdout),
		//pipe.Write(os.Stdout),
	)
	err := pipe.Run(p)

	// TODO fix everything below here
	if err != nil {
		return nil, err
	}
	return vid, nil
}

func chomp(line []byte) string {
	line = bytes.TrimRight(line, "\r\n")
	return strings.TrimSpace(string(line))
}

func Processor(f func(line []byte) []byte) pipe.Pipe {
	return pipe.TaskFunc(func(s *pipe.State) error {
		r := bufio.NewReader(s.Stdin)
		for {
			line, err := r.ReadBytes('\n')
			if len(line) > 0 {
				l := f(line)
				if len(l) > 0 {
					_, err := s.Stdout.Write(l)
					if err != nil {
						return err
					}
				}
			}
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
		}
	})
}

// A custom writer for debugging purposes
func CustomWrite(w io.Writer) pipe.Pipe {
	return pipe.TaskFunc(func(s *pipe.State) error {
		scanner := bufio.NewScanner(s.Stdin)
		for scanner.Scan() {
			_, err := w.Write([]byte("--->" + scanner.Text() + "\n"))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
