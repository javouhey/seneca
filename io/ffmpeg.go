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
    stdio "io"
    "fmt"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "reflect"
    "strings"
    "sync"
    "syscall"
    "time"

    "github.com/javouhey/seneca/util"
    "github.com/javouhey/seneca/vendor/launchpad.net/tomb"
)

var (
    ffmpegExec  string
    ffprobeExec string
)

const (
    CHUNK         = 1024
    APPDIR        = "seneca"
    PDIR          = "p"
    TMPMP4        = "temp.mp4"

    INVALID_VIDEO = "File %q not a recognizable video file\n\n%s\n"
    MISSING_PROG  = "Missing executable %q on your $PATH.\n\n%s\n"
)

func assignProgram(prog string, exec *string) {
    flag, err := util.IsExistProgram(prog)
    if err != nil {
        fmt.Fprintf(os.Stderr, MISSING_PROG, prog, util.ShortHelp)
        syscall.Exit(127)
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

// Dynamic values depending on time, inputs & OS
type Work struct {
    TmpDir  string
    PngDir  string
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
// @TODO allow only one time execution
func (v *VideoReader) Reset(size uint8) error {
    return v.reset2(size,
        func() string { return os.TempDir() },
        func() string { return string(os.PathSeparator) },
        func() int64 { t := time.Now(); return t.Unix() })
}

// compromise: no method overloading
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
    v.TmpDir = filepath.Join(tmpdir(), APPDIR, fmt.Sprintf("%d", (uniqnum())))
    v.PngDir = filepath.Join(v.TmpDir, PDIR)
    v.TmpFile = fmt.Sprintf("%s%0.2d%s", "img-%", size, "d.png")
    return nil
}

// Goal - Share memory by communicating
type FrameGenerator struct {}

// Task #1: Generate all the frames as PNGs
func (f FrameGenerator) Run(vr *VideoReader, args *util.Arguments) <-chan error {
    cmdFull := f.prepCli(vr, args)
    reply := make(chan error)
    go func() {
        if args.DryRun {
            fmt.Printf("  %s\n", cmdFull)
            reply <- nil
            return
        }

        if err := os.MkdirAll(vr.PngDir, os.ModePerm); err != nil {
            fmt.Fprintf(os.Stderr, "Unable to create %q\n\t%v\n", vr.PngDir, err)
            reply <- err
            return
        }

        cmd := exec.Command(ffmpegExec, cmdFull[1:]...)

        if err := cmd.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed executing %q\n\t%v\n", ffmpegExec, err)
            reply <- err
            return
        }
        if err := cmd.Wait(); err != nil {
            fmt.Fprintf(os.Stderr, "%q executed with errors\n\t%v\n", ffmpegExec, err)
            reply <- err
            return
        }
        reply <- nil
    }()
    return reply
}

func (f FrameGenerator) combineVf(args *util.Arguments) (bool, string) {
    needSpeed := !util.IsEmpty(args.SpeedSpec)
    var vfarg string
    switch {
    case args.NeedScaling:
        vfarg += args.ScaleFilter
        if needSpeed {
            vfarg += "," + args.SpeedSpec
        }
        return true, vfarg
    case needSpeed:
        vfarg += args.SpeedSpec
        return true, vfarg
    }
    return false, ""
}

func (f FrameGenerator) prepCli(vr *VideoReader, args *util.Arguments) []string {
    cmdFull := []string{ffmpegExec, "-ss", args.From.String()}

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

    cmdFull = append(cmdFull, "-i", vr.Filename, "-an")

    if vf, s := f.combineVf(args); vf {
        cmdFull = append(cmdFull, "-vf", s)
    }

    cmdFull = append(cmdFull, "-q:v", "2", "-f", "image2", "-vsync", "cfr")
    cmdFull = append(cmdFull, "-r", fmt.Sprintf("%d", args.Fps), "-y")
    cmdFull = append(cmdFull, "-progress", fmt.Sprintf("http://127.0.0.1:%d",
                                                       args.Port))
    vr.Reset(uint8(f.guess(secs)))
    cmdFull = append(cmdFull, filepath.Join(vr.PngDir, vr.TmpFile))

    if args.Verbose {
        fmt.Printf("  Workdir: %q\n", vr.TmpDir)
        fmt.Printf("   Frames: %q\n", vr.TmpFile)
        fmt.Printf("      gif: %q\n", vr.Gif)
    }
    return cmdFull
}

// Naive way to guess how many images are
// captured per frames.
func (f FrameGenerator) guess(secs float64) int {
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

// Goal - Communicate by sharing memory
type Muxer struct {
    err error
    sync.Mutex
}

func (m Muxer) prepCli(vr *VideoReader, args *util.Arguments) []string {
    cmdFull := []string{ffmpegExec, "-f", "image2", "-y"}
    cmdFull = append(cmdFull, "-progress", fmt.Sprintf("http://127.0.0.1:%d",
                                                       args.Port))
    cmdFull = append(cmdFull, "-i", filepath.Join(vr.PngDir, vr.TmpFile))
    cmdFull = append(cmdFull, "-c:v", "libx264", "-crf", "23")
    cmdFull = append(cmdFull, "-vf",
        fmt.Sprintf("fps=%d,format=yuv420p", args.Fps))
    cmdFull = append(cmdFull, "-preset", "veryslow")
    cmdFull = append(cmdFull, filepath.Join(vr.TmpDir, TMPMP4))
    return cmdFull
}

func (m *Muxer) setError(e error) {
    m.Lock()
    defer m.Unlock()
    m.err = e
}

// NOTE: receiver doesn't need to be a reference
func (m *Muxer) Error() error {
    m.Lock()
    defer m.Unlock()
    return m.err
}

// a priori: FrameGenerator task was executed without errors
func (m *Muxer) Run(vr *VideoReader, args *util.Arguments) *sync.WaitGroup {
    var wg sync.WaitGroup

    cmdFull := m.prepCli(vr, args)
    wg.Add(1)

    go func() {
        defer wg.Done()
        if args.DryRun {
            fmt.Printf("  %s\n", cmdFull)
            return
        }

        cmd := exec.Command(ffmpegExec, cmdFull[1:]...)

        if err := cmd.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed executing %q\n\t%v\n", ffmpegExec, err)
            m.setError(err)
            return
        }
        if err := cmd.Wait(); err != nil {
            fmt.Fprintf(os.Stderr, "%q executed with errors\n\t%v\n", ffmpegExec, err)
            m.setError(err)
            return
        }
    }()
    return &wg
}

type GifWriter struct {
    Tombstone tomb.Tomb
}

func (g *GifWriter) Stop() error {
    g.Tombstone.Kill(nil)
    return g.Tombstone.Wait()
}

func (g *GifWriter) Run(vr *VideoReader, args *util.Arguments) {
    cmdFull := g.prepCli(vr, args)

    go func() {
        defer g.Tombstone.Done()

        if args.DryRun {
            fmt.Printf("  %s\n", cmdFull)
            g.Tombstone.Kill(nil)
            return
        }

        time.Sleep(2 * time.Second)

        // Cooperative cancelation.
        select {
        case <- g.Tombstone.Dying():
            fmt.Println("aborting")
            return
        default:
            // noop
        }

        cmd := exec.Command(ffmpegExec, cmdFull[1:]...)

        if err := cmd.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed executing %q\n\t%v\n", ffmpegExec, err)
            g.Tombstone.Kill(err)
            return
        }
        if err := cmd.Wait(); err != nil {
            fmt.Fprintf(os.Stderr, "%q executed with errors\n\t%v\n", ffmpegExec, err)
            g.Tombstone.Kill(err)
            return
        }
    }()
}

func (g GifWriter) prepCli(vr *VideoReader, args *util.Arguments) []string {
    cmdFull := []string{ffmpegExec, "-i"}
    cmdFull = append(cmdFull, filepath.Join(vr.TmpDir, TMPMP4))
    cmdFull = append(cmdFull, "-progress", fmt.Sprintf("http://127.0.0.1:%d",
                                                       args.Port))
    cmdFull = append(cmdFull, "-y", "-vf", "format=rgb24")
    cmdFull = append(cmdFull, filepath.Join(vr.TmpDir, vr.Gif))
    return cmdFull
}

// getMetadata parses output of `ffprobe` into a VideoReader
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
