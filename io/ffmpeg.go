package io

import (
    "time"
    "fmt"
    "os/exec"
    "io"
    "log"

    "github.com/javouhey/seneca/util"
)

var (
    ffmpegExec = ""
    ffprobeExec = ""
)

const (
    A = 1
)

func assignProgram(prog string) {
    flag, err := util.IsExistProgram(prog) 
    if err != nil {
        log.Fatalf("Cannot find program '%#v' on your system. %#v", prog, err)
    }
    if flag {
        ffprobeExec = prog
    }
}

func init() {
    fmt.Println("init() from ffmpeg.go", time.Now())
    assignProgram("ffprobe")
    fmt.Printf("%#v\n", ffprobeExec)
    assignProgram("ffmpeg")
    fmt.Printf("%#v\n", ffprobeExec)
}

type VideoSize struct {
  width, height uint8
}

type VideoReader struct {
  filename string
  fps float64
  numberofframes uint16
  duration float64
  VideoSize
}



func test(filename string) error {
  cmd := exec.Command("ffmpeg", "-i", filename)
  //cmd := exec.Command("ls", "-l", "opop")
  fmt.Println("test")
  stderr, _ := cmd.StderrPipe()
  err := cmd.Start()
  if err != nil {
    fmt.Println("errrrrrr")
  }
  var b []byte
  b = make([]byte, 4196)
  //go io.Copy(os.Stderr, stderr)
  go io.ReadFull(stderr, b)
  cmd.Wait()
  fmt.Println("res", "e", string(b))
  return nil
}

func New(filename string) (*VideoReader, error) {
  test(filename)
  r := new(VideoReader)
  return r, nil
}
