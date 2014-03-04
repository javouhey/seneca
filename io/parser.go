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
  "fmt"
  "io"
  //"log"
  "os"
  "regexp"
  "strings"

  "labix.org/v2/pipe"
)

const (
  sVideo    = "Video:"
  sDuration = "Duration:"
)

var (
  duration = regexp.MustCompile(`Duration: (?P<duration>\d{2}:\d{2}:\d{2}.\d{2})`)
)

// Returns only the items that we need
func parse(str string) map[string]string {
  scanner := bufio.NewScanner(strings.NewReader(str))
  for scanner.Scan() {
    fmt.Println("=> ", scanner.Text())
  }
  return nil
}

func parse2(data *bytes.Buffer) (VideoSize, error) {
  fmt.Println("parse2")

  p := pipe.Line(
    pipe.Read(bytes.NewReader(data.Bytes())),

    pipe.Filter(func(line []byte) bool {
      s := string(line)
      s = strings.TrimSpace(s)
      return strings.HasPrefix(s, sDuration) ||
             strings.Index(s, sVideo) >= 0
    }),

    CustomWrite(os.Stdout),
    //pipe.Write(os.Stdout),
  )
  err := pipe.Run(p)
  if err != nil {
    return VideoSize{}, err
  }

  return VideoSize{1, 2}, nil
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
