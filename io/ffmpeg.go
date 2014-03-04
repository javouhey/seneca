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
  "bytes"
  "fmt"
  "github.com/javouhey/seneca/util"
  stdio "io"
  "log"
  "os/exec"
  "reflect"
  "time"
)

var (
  ffmpegExec  = ""
  ffprobeExec = ""
)

const (
  CHUNK = 1024
)

func assignProgram(prog string, exec *string) {
  flag, err := util.IsExistProgram(prog)
  if err != nil {
    log.Fatalf("Cannot find program '%#v' on your system. %#v", prog, err)
  }
  if flag {
    p := reflect.ValueOf(exec)
    p.Elem().SetString(prog)
  }
}

func init() {
  fmt.Println("init() from ffmpeg.go", time.Now())
  assignProgram("ffprobe", &ffprobeExec)
  assignProgram("ffmpeg", &ffmpegExec)
}

type VideoSize struct {
  width, height uint8
}

type VideoReader struct {
  filename       string
  fps            float64
  numberofframes uint16
  duration       float64
  VideoSize
}

// getMetadata parses output of `ffprobe` into a map
func getMetadata(videoFile string) (map[string]string, error) {
  var n int
  var err error
  var stderr stdio.ReadCloser
  result := make(map[string]string)

  cmd := exec.Command(ffprobeExec, videoFile)
  if stderr, err = cmd.StderrPipe(); err != nil {
    return result, err
  }

  if err := cmd.Start(); err != nil {
    return result, err
  }

  var data bytes.Buffer
  for {
    n = 0
    err = nil
    tmp := make([]byte, CHUNK)
    n, err = stdio.ReadFull(stderr, tmp)
    if err == nil {
      data.Write(tmp)
    } else {
      if err == stdio.ErrUnexpectedEOF {
        if n > 0 {
          tmp = tmp[0:n]
          data.Write(tmp)
        }
        break
      }
    }
  }
  //fmt.Printf("Buffer size %d\n", data.Len())
  cmd.Wait()

  // TODO: finish this!
  //parse(data.String())

  parse2(&data)

  return result, nil
}

func NewVideoReader(filename string) (*VideoReader, error) {
  getMetadata(filename)
  r := new(VideoReader)
  return r, nil
}
