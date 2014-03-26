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

    if len(os.Args) == 1 {
        fmt.Printf("%s", util.HelpMessage)
        syscall.Exit(0)
    }

    args := util.NewArguments()
    if err := args.Parse(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n\n%s\n", err, util.ShortHelp)
        syscall.Exit(1)
    }

    fmt.Printf("%#v\n", args)

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
        syscall.Exit(2)
    }

    var vr *io.VideoReader
    var errVr error

    filename, _ := util.SanitizeFile(args.VideoIn)
    vr, errVr = io.NewVideoReader(filename)
    if errVr != nil {
        fmt.Fprintf(os.Stderr, io.INVALID_VIDEO, filename, util.ShortHelp)
        syscall.Exit(128)
    }

    ValidateWithVideo(vr, args)
    fmt.Printf("%#v\n", vr)

    // --- setup our progress bar ---
    listener := NewListener(args.Port)

    defer func() {
        close(ipc)
        log.Printf("Closed progress channel")
        listener.Close()
        log.Printf("Closed TCP listener")
    }()

    go progress.Outputter(ipc)
    go progress.Progress(listener, ipc, args.Port)

    // --- ffmpeg execution ---

    fmt.Println("@TODO -- start goroutine to execute ffmpeg")

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

// TODO 1. Ensure that the provided -from and -length does not exceed
//         the duration of this video.
func ValidateWithVideo(vr *io.VideoReader, args *util.Arguments) error {
    fmt.Println("ValidateWithVideo")
    return nil
}
