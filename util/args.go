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

package util

import (
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    "time"
)

type Instant []time.Duration

func (i *Instant) String() string {
    return fmt.Sprint(*i)
}

func (i *Instant) Set(value string) error {
    if len(*i) > 0 {
        return errors.New("interval flag already set")
    }
    for _, dt := range strings.Split(value, ",") {
        duration, err := time.ParseDuration(dt)
        if err != nil {
            return err
        }
        *i = append(*i, duration)
    }
    return nil
}

type Arguments struct {
    Help     bool
    Version  bool
    VideoIn  string
    Port     int
    FromTime Instant
}

func New() *Arguments {
    args := new(Arguments)
    return args
}

func (a *Arguments) Parse(arguments []string) error {
    f := flag.NewFlagSet("seneca", flag.ContinueOnError)
    f.SetOutput(ioutil.Discard)

    f.BoolVar(&a.Help, "h", false, "")
    f.BoolVar(&a.Version, "version", false, "")
    f.StringVar(&a.VideoIn, "video-infile", a.VideoIn, "")
    f.IntVar(&a.Port, "port", 8080, "")
    f.Var(&a.FromTime, "from", "")

    if err := f.Parse(arguments); err != nil {
        return err
    }

    return nil
}

func (a *Arguments) Validate() error {
    // TODO
    return nil
}
