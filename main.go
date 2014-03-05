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
  "log"
  "flag"
  "fmt"
  "net"
  "syscall"
  //"sync"
  //"time"

  "github.com/javouhey/seneca/io"
  "github.com/javouhey/seneca/util"
  "github.com/javouhey/seneca/progress"
)

var (
    GitSHA string
    Version string
)

func main() {

  var (
    flVersion = flag.Bool("ver", false, "Print version")
    flVideo = flag.String("video", "", "relative or full path to the video file")
  )
  flag.Parse()

  if *flVersion {
    printVersion()
    syscall.Exit(0)
  }

  filename, err := util.SanitizeFile(*flVideo)
  if err != nil {
    log.Fatalf("The video file provided is invalid (%s)", err.Error())
  }


// --- ensure file is valid video --

  vr, err2 := io.NewVideoReader(filename)
  if err2 != nil {
    log.Fatalf("Not a video file (%s)", err2.Error())
  }
  fmt.Printf("%#v\n", vr)

// --- setup our progress bar ---

  mychan := make(chan progress.Status)
  l, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Fatal(err)
  }

  defer func() {
    close(mychan)
    log.Printf("Closed progress channel")
    l.Close()
    log.Printf("Closed TCP listener")
  }()

  go progress.Outputter(mychan)
  go progress.Progress(l, mychan)

// --- block ---

  var input string
  fmt.Scanln(&input)
}

func init() {
  //fmt.Println("init() from main.go", time.Now())
}

func printVersion() {
  fmt.Printf("Seneca version %s, git SHA %s\n", Version, GitSHA)
}
