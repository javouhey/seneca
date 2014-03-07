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
    "flag"
    "fmt"
    "log"
    "net"
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
)

func main() {

    var (
        flVersion = flag.Bool("ver", false, "Print version")
        flVideo   = flag.String("video", "", "relative or full path to the video file")
        flPort    = flag.Int("port", 8080, "http listener port for progress reporter")
        flStart   = flag.String("start", "00:00:00", "instant within video to start capture")
    )
    flag.Parse()

    if *flVersion {
        printVersion()
        syscall.Exit(0)
    }

    // --- ensure input video is valid file #1 --
    filename, err := util.SanitizeFile(*flVideo)
    if err != nil {
        log.Fatalf("The video file provided is invalid (%s)", err.Error())
    }

    // --- valid HTTP port ---
    util.ValidatePort(*flPort)

    // --- ensure input video is valid file #2 --
    vr, err2 := io.NewVideoReader(filename)
    if err2 != nil {
        log.Fatalf("Not a video file (%s)", err2.Error())
    }
    fmt.Printf("%#v\n", vr)

    // --- valid start time ---
    util.ParseStartTime(*flStart, vr.Duration)

    // --- setup our progress bar ---
    listener := NewListener(*flPort)

    defer func() {
        close(ipc)
        log.Printf("Closed progress channel")
        listener.Close()
        log.Printf("Closed TCP listener")
    }()

    go progress.Outputter(ipc)
    go progress.Progress(listener, ipc, *flPort)

    // --- block wait ---
    var input string
    fmt.Scanln(&input)
}

func init() {
    ipc = make(chan progress.Status)
}

func NewListener(port int) net.Listener {
    listener, err := net.Listen("tcp", util.ToPort(port))
    if err != nil {
        log.Fatal(err)
    }
    return listener
}

func printVersion() {
    fmt.Printf("Seneca version %s, git SHA %s\n", Version, GitSHA)
}
