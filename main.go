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

package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "path/filepath"
    "runtime"
    "syscall"

    "github.com/javouhey/seneca/io"
    "github.com/javouhey/seneca/progress"
    "github.com/javouhey/seneca/util"
)

var (
    GitSHA  string
    Version string

    listener net.Listener
    ipc      chan progress.Status

    task1    *io.FrameGenerator
    task2    *io.Muxer
    task3    *io.GifWriter
)

func main() {

    if len(os.Args) == 1 {
        fmt.Printf("%s", util.HelpMessage)
        syscall.Exit(0)
    }

    args := util.NewArguments()
    if err := args.Parse(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n\n%s\n", err, util.ShortHelp)
        syscall.Exit(1)
    }

    if args.Verbose {
        fmt.Printf("  %#v\n", args)
    }

    if args.Version {
        printVersion()
        syscall.Exit(0)
    }

    if args.Help {
        fmt.Printf("%s", util.HelpMessage)
        syscall.Exit(0)
    }

    if err := args.Validate(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n\n%s\n", err, util.ShortHelp)
        syscall.Exit(1)
    }

    var vr *io.VideoReader
    var errVr error

    filename, _ := util.SanitizeFile(args.VideoIn)
    vr, errVr = io.NewVideoReader(filename, args.DryRun)
    if errVr != nil {
        fmt.Fprintf(os.Stderr, io.INVALID_VIDEO, filename, util.ShortHelp)
        syscall.Exit(1)
    }

    if args.Verbose {
        fmt.Printf("%s", vr)
    }

    // --- setup progress notification ---
    listener := NewTCPListener(args.Port)

    defer func() {
        cleanup(vr)

        listener.Close()
        if args.Verbose {
            fmt.Println("Closed TCP listener")
        }
        close(ipc)
        if args.Verbose {
            fmt.Println("Closed progress channel")
        }
    }()

    go progress.StatusLogger(ipc)
    go progress.Progress(listener, ipc, args.Port)

    // --- Pipeline ---
    reply := task1.Run(vr, args)
    if err := <- reply; err != nil {
        syscall.Exit(126)
    }

    wg := task2.Run(vr, args)
    wg.Wait()
    if err := task2.Error(); err != nil {
        syscall.Exit(126)
    }

    task3.Run(vr, args)
    /* 
    Sample code for cancelling a goroutine

    time.Sleep(1 * time.Second)
    log.Fatal(task3.Stop())
    */
    if err := task3.Tombstone.Wait(); err != nil {
        syscall.Exit(126)
    }

    sayGoodbye(vr)
}

func init() {
    ipc = make(chan progress.Status)
    task1 = new(io.FrameGenerator)
    task2 = new(io.Muxer)
    task3 = new(io.GifWriter)
    runtime.GOMAXPROCS(3)
}

func NewTCPListener(port int) net.Listener {
    listener, err := net.Listen("tcp", util.ToPort(port))
    if err != nil {
        log.Fatal(err)
    }
    return listener
}

func cleanup(vr *io.VideoReader) {
    if vr != nil && !util.IsEmpty(vr.PngDir) {
        if err := os.RemoveAll(vr.PngDir); err != nil {
            fmt.Printf("WARNING: Removing %s encountered errors\n", vr.PngDir)
            fmt.Printf("\t%s", err.Error())
        }
    }
}

func sayGoodbye(vr *io.VideoReader) {
    if vr != nil && !util.IsEmpty(vr.TmpDir) {
        fmt.Println("\n\nYour animated GIF is ready at location:")
        fmt.Printf("  %s\n\n", filepath.Join(vr.TmpDir, vr.Gif))
    }
}

func printVersion() {
    fmt.Printf("\nSeneca version %s, git SHA %s\n", Version, GitSHA)
}

//@TODO 1. Ensure that the provided -from and -length does not exceed
//         the duration of this video.
//func ValidateWithVideo(vr *io.VideoReader, args *util.Arguments) error {
//    return nil
//}
