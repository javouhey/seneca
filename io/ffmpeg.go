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
    "github.com/javouhey/seneca/util"
    stdio "io"
    //"log"
    "fmt"
    "os"
    "os/exec"
    "reflect"
    "syscall"
    "time"
)

var (
    ffmpegExec  string
    ffprobeExec string
)

const (
    CHUNK         = 1024
    INVALID_VIDEO = "File %q not a recognizable video file\n\n%s\n"
    MISSING_PROG  = "Missing executable %q on your $PATH.\n\n%s\n"
)

func assignProgram(prog string, exec *string) {
    flag, err := util.IsExistProgram(prog)
    if err != nil {
        fmt.Fprintf(os.Stderr, MISSING_PROG, prog, util.ShortHelp)
        syscall.Exit(128)
    }
    if flag {
        p := reflect.ValueOf(exec)
        p.Elem().SetString(prog)
    }
}

func init() {
    assignProgram("ffprobe", &ffprobeExec)
    assignProgram("ffmpeg", &ffmpegExec)
}

type VideoSize struct {
    Width, Height uint16
}

type Work struct {
    TmpDir  string
    TmpFile string
}

type VideoReader struct {
    Filename string
    Fps      float32
    Duration time.Duration
    VideoSize
    Work
}

// Generates internally the temporary work directories
// and other runtime constants etc.
func (v *VideoReader) Reset(size uint8) {
    v.reset2(size,
        func() string { return os.TempDir() },
        func() string { return string(os.PathSeparator) },
        func() int64 { t := time.Now(); return t.Unix() })
}

func (v *VideoReader) reset2(size uint8,
    tmpdir func() string,
    pathsep func() string,
    uniqnum func() int64) {
    v.TmpDir = fmt.Sprintf("%s%s%d", tmpdir(),
        pathsep(), uniqnum())
    v.TmpFile = fmt.Sprintf("%s%0.2d%s", "img-%", size, "d.png")
}

// generate all the frames as PNGs
// run as a goroutine
func GenerateFrames(vr *VideoReader) {
    // assemble all the command arguments
}

// getMetadata parses output of `ffprobe` into a map
func getMetadata(videoFile string) (*VideoReader, error) {
    cmd := exec.Command(ffprobeExec, videoFile)
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return nil, err
    }

    if err := cmd.Start(); err != nil {
        return nil, err
    }

    var (
        data bytes.Buffer
        n    int
    )
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

    if err := cmd.Wait(); err != nil {
        return nil, err
    }

    vr, err := parse(&data)
    return vr, err
}

func NewVideoReader(filename string) (vr *VideoReader, err error) {
    vr, err = getMetadata(filename)
    if err == nil {
        vr.Filename = filename
    }
    return
}
