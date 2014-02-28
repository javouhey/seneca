package main

import (
  "time"
  "fmt"
  "github.com/javouhey/animfilereader/io"
)

func main() {
  //io.New("budapest.flv")
  io.New("plank.mp4")
}

func init() {
  fmt.Println("init() from main.go", time.Now())
}
