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
  //"log"
  "os"
  "regexp"
  "strings"
  "time"

  "github.com/javouhey/seneca/vendor/labix.org/v2/pipe"
  "github.com/javouhey/seneca/util"
)

const (
  sVideo    = "Video:"
  sDuration = "Duration:"
)

var (
  // Duration: 00:08:20.53, start: 0.000000, bitrate: 709 kb/s
  Regex1 = regexp.MustCompile(`^Duration: (?P<duration>\d{2}:\d{2}:\d{2}).\d{2}(.*)$`)

  InvalidDuration = errors.New("Duration input is invalid")
)

// Returns only the items that we need
func parse(str string) map[string]string {
  scanner := bufio.NewScanner(strings.NewReader(str))
  for scanner.Scan() {
    fmt.Println("=> ", scanner.Text())
  }
  return nil
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

  retval, err := time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
  if err != nil {
    return 0, err
  }
  return retval, nil
}

func parse2(data *bytes.Buffer) (VideoSize, error) {
  fmt.Println("parse2")

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
      line = bytes.TrimRight(line, "\r\n")
      s := string(line)
      s = strings.TrimSpace(s)
      if strings.HasPrefix(s, sDuration) {
        d, _ := ParseDuration(s)
        // TODO log the err
        vid.duration = d
        return make([]byte, 0)
      } else {
        return line // leave untouched
      }
    }),

    CustomWrite(os.Stdout),
    //pipe.Write(os.Stdout),
  )
  err := pipe.Run(p)
  if err != nil {
    return VideoSize{}, err
  }
  fmt.Println("vid", vid)
  return VideoSize{1, 2}, nil
}


func Processor(f func(line []byte) []byte) pipe.Pipe {
  return pipe.TaskFunc(func(s *pipe.State) error {
    r := bufio.NewReader(s.Stdin)
    for {
      line, err := r.ReadBytes('\n')
      if len(line) > 0 {
        line := f(line)
        if len(line) > 0 {
          _, err := s.Stdout.Write(line)
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
    return nil
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
