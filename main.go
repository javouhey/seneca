package main

import (
  "time"
  "fmt"
  "github.com/javouhey/seneca/util"
  "github.com/javouhey/seneca/io"
)

func main() {
  //io.New("budapest.flv")
  util.IsEmpty("")
  io.New("plank.mp4")
}

func init() {
  fmt.Println("init() from main.go", time.Now())
}
