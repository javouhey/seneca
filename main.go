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
  "time"
  "fmt"
  "github.com/javouhey/seneca/io"
)

func main() {
  fmt.Printf("%#v\n", simplyExec("goproplane.mp4"))
  fmt.Printf("%#v\n", simplyExec("budapest.flv"))
  fmt.Printf("%#v\n", simplyExec("plank.mp4"))
}

func simplyExec(filename string) *io.VideoReader {
  res, err := io.NewVideoReader(filename)
  if err != nil {
    panic ("fail to get reader: " + err.Error())
  }
  return res
}

func init() {
  fmt.Println("init() from main.go", time.Now())
}
