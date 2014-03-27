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
    "path"
    "reflect"
    "strings"
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

// Dynamic values depending on time, & OS.
type Work struct {
    TmpDir  string
    TmpFile string
    Gif     string
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
func (v *VideoReader) Reset(size uint8) error {
    return v.reset2(size,
        func() string { return os.TempDir() },
        func() string { return string(os.PathSeparator) },
        func() int64 { t := time.Now(); return t.Unix() })
}

func (v *VideoReader) reset2(size uint8,
    tmpdir func() string,
    pathsep func() string,
    uniqnum func() int64) error {

    if util.IsEmpty(v.Filename) {
        return fmt.Errorf("Missing VideoReader.Filename")
    }
    _, video := path.Split(v.Filename)
    parts := strings.Split(video, ".")
    if len(parts) < 2 {
        return fmt.Errorf("Invalid VideoReader.Filename")
    }
    name := strings.TrimSpace(parts[0])
    if util.IsEmpty(name) {
        return fmt.Errorf("Empty VideoReader.Filename")
    }
    v.Gif = name + ".gif"
    v.TmpDir = fmt.Sprintf("%s%s%d", tmpdir(), pathsep(), uniqnum())
    v.TmpFile = fmt.Sprintf("%s%0.2d%s", "img-%", size, "d.png")
    return nil
}


// generate all the frames as PNGs
// run as a goroutine
func GenerateFrames(vr *VideoReader, args *util.Arguments) {
    cmdFull := []string{ffmpegExec, "-i", vr.Filename, "-an"}
    if args.NeedScaling {
        cmdFull = append(cmdFull, "-vf", args.ScaleFilter)
    }
    cmdFull = append(cmdFull, "-ss", args.From.String())

    secs := args.Length.Seconds()
    switch {
    case secs < 60.0 && secs > 0.0:
        cmdFull = append(cmdFull, "-t", fmt.Sprintf("%d", int64(secs)))
    case args.Length > vr.Duration:
        fallthrough
    default:
        fmt.Fprintf(os.Stderr, "WARNING: %d secs is outside of range. " +
                               "Forcing to 3 secs.\n", int64(secs))
        cmdFull = append(cmdFull, "-t", "3")
    }
    cmdFull = append(cmdFull, "-q:v", "2", "-f", "image2", "-vsync", "cfr")
    cmdFull = append(cmdFull, "-r", fmt.Sprintf("%d", args.Fps), "-y")
    cmdFull = append(cmdFull, "-progress", fmt.Sprintf("http://127.0.0.1:%d",
                                                       args.Port))
    vr.Reset(uint8(guess(secs)))
    cmdFull = append(cmdFull,
        fmt.Sprintf("%s%s%s",
            vr.TmpDir,
            string(os.PathSeparator),
            vr.TmpFile))

    if args.DryRun {
        fmt.Printf("  %s\n", cmdFull)
        return
    }

    if args.Verbose {
        fmt.Printf("  Workdir: %q\n", vr.TmpDir)
        fmt.Printf("   Frames: %q\n", vr.TmpFile)
        fmt.Printf("      gif: %q\n", vr.Gif)
    }

    if err := os.MkdirAll(vr.TmpDir, os.ModePerm); err != nil {
        fmt.Printf("Mkdir? %q\n", err.Error())
        return
    }

    cmd := exec.Command(ffmpegExec, cmdFull[1:]...)
    if err := cmd.Start(); err != nil {
        fmt.Printf("Exec? %q\n", err.Error())
    }
    if err := cmd.Wait(); err != nil {
        fmt.Printf("%q\n", err.Error())
    }

}

func guess(secs float64) int {
    switch {
    case secs > 0.0 && secs < 15.0:
        return 3
    case secs >= 15.0 && secs < 30.0:
        return 4
    case secs >= 30.0 && secs < 60.0:
        fallthrough
    default:
        return 5
    }
}

func MergeAsVideo(vr *VideoReader, args *util.Arguments) {
    cmdFull := []string{ffmpegExec, "-f", "image2", "-y"}
    cmdFull = append(cmdFull, "-progress", fmt.Sprintf("http://127.0.0.1:%d",
                                                       args.Port))
    cmdFull = append(cmdFull, "-i",
        fmt.Sprintf("%s%s%s",
            vr.TmpDir,
            string(os.PathSeparator),
            vr.TmpFile))

    cmdFull = append(cmdFull, "-c:v", "libx264", "-crf", "23")
    cmdFull = append(cmdFull, "-vf",
        fmt.Sprintf("fps=%d,format=yuv420p", args.Fps))
    cmdFull = append(cmdFull, "-preset", "veryslow")
    cmdFull = append(cmdFull,
        fmt.Sprintf("%s%s%s",
            vr.TmpDir,
            string(os.PathSeparator),
            "temp.mp4"))
    /*
    -i /tmp/mar20-zlatan/x/img-%03d.png 
    -c:v libx264  -crf 23 
    -vf "fps=17,format=yuv420p" -preset veryslow  zlatan1.mp4 
    */
    if args.DryRun {
        fmt.Printf("  %s\n", cmdFull)
        return
    }
    cmd := exec.Command(ffmpegExec, cmdFull[1:]...)
    if err := cmd.Start(); err != nil {
        fmt.Printf("Exec? %q\n", err.Error())
    }
    if err := cmd.Wait(); err != nil {
        fmt.Printf("%q\n", err.Error())
    }
}

// getMetadata parses output of `ffprobe` into a map
func getMetadata(videoFile string, dryRun bool) (*VideoReader, error) {
    cmdFull := []string{ffprobeExec, videoFile}
    if dryRun {
        fmt.Printf("  %s\n", cmdFull)
    }
    cmd := exec.Command(ffprobeExec, cmdFull[1:]...)
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

func NewVideoReader(filename string, dryRun bool) (vr *VideoReader, err error) {
    vr, err = getMetadata(filename, dryRun)
    if err == nil {
        vr.Filename = filename
    }
    return
}
