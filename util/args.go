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

// Argument parsing and utility functions
package util

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var empty = struct{}{}

type Arguments struct {
	Help    bool
	Version bool
	DryRun  bool
	Verbose bool
	VideoIn string
	Port    int

	NeedScaling bool
	ScaleFilter string
	Fps         int
	SpeedSpec   string

	From   TimeCode
	Length time.Duration
}

func NewArguments() *Arguments {
	args := new(Arguments)
	return args
}

func (a *Arguments) Parse(arguments []string) error {
	f := flag.NewFlagSet("seneca", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	f.BoolVar(&a.Help, "h", false, "")
	f.BoolVar(&a.Version, "version", false, "")
	f.BoolVar(&a.DryRun, "dry-run", false, "")
	f.BoolVar(&a.Verbose, "vv", false, "")
	f.StringVar(&a.VideoIn, "video-infile", a.VideoIn, "")
	f.IntVar(&a.Port, "port", 8080, "")

	scalingArg := f.String("scale", "_:_", "")
	speedArg := f.String("speed", "placebo", "")
	f.IntVar(&a.Fps, "fps", 25, "")

	f.DurationVar(&a.Length, "length", 3*time.Second, "")
	fromArg := f.String("from", "00:00:00", "")

	if err := f.Parse(arguments); err != nil {
		return err
	}

	if err := preprocessScale(a, *scalingArg); err != nil {
		return err
	}
	if err := preprocessSpeed(a, *speedArg); err != nil {
		return err
	}
	if err := preprocessFrom(a, *fromArg); err != nil {
		return err
	}

	return nil
}

func (a *Arguments) Validate() error {
	if _, err := SanitizeFile(a.VideoIn); err != nil {
		return err
	}

	if err := ValidatePort(a.Port); err != nil {
		return err
	}

	if a.Fps < 1 || a.Fps > 30 {
		return fmt.Errorf("frame rate -fps %d not in range [1, 30]", a.Fps)
	}

	return nil
}

func preprocessFrom(a *Arguments, fromArg string) error {
	if fromArg != "00:00:00" {
		tc, err := ParseFrom(fromArg)
		if err != nil {
			return err
		}
		a.From = *tc
	}
	return nil
}

var rgxScale = regexp.MustCompile(`^(?P<width>(_|\d{1,})):(?P<height>(_|\d{1,}))$`)

var isUnderscore = func(arg string) bool {
	return "_" == arg
}

func preprocessSpeed(a *Arguments, speedArg string) error {
	if speedArg != "placebo" {
		vf, err := DecodeSpeed(speedArg)
		if err != nil {
			return err
		}
		a.SpeedSpec = vf
	}
	return nil
}

func preprocessScale(a *Arguments, scalingArg string) error {
	err := fmt.Errorf("BAD arg to -scale %q", scalingArg)
	if scalingArg != "_:_" {
		if !rgxScale.MatchString(scalingArg) {
			return err
		}
		w := rgxScale.ReplaceAllString(scalingArg,
			fmt.Sprintf("${%s}", rgxScale.SubexpNames()[1]))
		h := rgxScale.ReplaceAllString(scalingArg,
			fmt.Sprintf("${%s}", rgxScale.SubexpNames()[3]))

		var (
			v1, v2 uint64
			vf     string
		)
		switch {
		case isUnderscore(w) && !isUnderscore(h):
			if v1, err = strconv.ParseUint(h, 10, 16); err != nil {
				return fmt.Errorf("height in -scale %q overflow",
					scalingArg)
			}
			if vf, err = HeightOnly.Decode(uint16(v1)); err != nil {
				return err
			}
			a.ScaleFilter = vf

		case !isUnderscore(w) && isUnderscore(h):
			if v2, err = strconv.ParseUint(w, 10, 16); err != nil {
				return fmt.Errorf("width in -scale %q overflow",
					scalingArg)
			}
			if vf, err = WidthOnly.Decode(uint16(v2)); err != nil {
				return err
			}
			a.ScaleFilter = vf

		default:
			if v1, err = strconv.ParseUint(h, 10, 16); err != nil {
				return fmt.Errorf("height in -scale %q overflow",
					scalingArg)
			}
			if v2, err = strconv.ParseUint(w, 10, 16); err != nil {
				return fmt.Errorf("width in -scale %q overflow",
					scalingArg)
			}
			if vf, err = WidthHeight.Decode(uint16(v2), uint16(v1)); err != nil {
				return err
			}
			a.ScaleFilter = vf
		}

		a.NeedScaling = true
	}
	return nil
}

/////////////////////////////////////////////////////////////////

// @deprecated
type Instant []time.Duration

// @deprecated
func (i *Instant) String() string {
	return fmt.Sprint(*i)
}

// @deprecated
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

/////////////////////////////////////////////////////////////////

const (
	tclayout = "15:04:05"
)

type TimeCode time.Time

func (tc TimeCode) String() string {
	t := time.Time(tc)
	return fmt.Sprintf("%0.2d:%0.2d:%0.2d",
		t.Hour(), t.Minute(), t.Second())
}

func ParseFrom(arg string) (*TimeCode, error) {
	t, err := time.Parse(tclayout, arg)
	if err != nil {
		return nil, err
	}
	result := TimeCode(t)
	return &result, nil
}

/////////////////////////////////////////////////////////////////

var speeds = map[string]struct{}{
	"veryfast": empty,
	"faster":   empty,
	"placebo":  empty,
	"slower":   empty,
	"veryslow": empty,
}

//var ErrInvalidSpeed = errors.New("Speed not recognized")

// Converts speed specification to ffmpeg option
//    speedup: -vf "setpts=(1/X)*PTS"
//   slowdown: -vf "setpts=(X/1)*PTS"
func DecodeSpeed(arg string) (string, error) {
	if _, ok := speeds[arg]; !ok {
		return "", fmt.Errorf("Invalid speed argument %q", arg)
	}

	switch arg {
	case "veryfast":
		return "setpts=1/3*PTS", nil
	case "faster":
		return "setpts=1/2*PTS", nil
	case "placebo":
		return "", nil
	case "slower":
		return "setpts=2*PTS", nil
	case "veryslow":
		return "setpts=3*PTS", nil
	}

	return "", fmt.Errorf("Invalid speed argument %q", arg)
}

/////////////////////////////////////////////////////////////////

type predicate func(uint16) bool

type ScaleType uint16

const (
	// -scale width:_
	//   height scaled to keep aspect ratio
	WidthOnly ScaleType = 1 << iota

	// -scale _:height
	//   Width scaled to keep aspect ratio
	HeightOnly

	// -scale width:height
	WidthHeight
)

var scales = map[ScaleType]struct{}{
	WidthOnly:   empty,
	HeightOnly:  empty,
	WidthHeight: empty,
}

// Converts into a valid argument to the -vf option of ffmpeg
func (s ScaleType) interpolate(width, height uint16) string {
	switch s {
	case WidthHeight:
		return fmt.Sprintf("scale=%d:%d", width, height)
	case HeightOnly:
		return fmt.Sprintf("scale=trunc(oh*a/2)*2:%d", height)
	case WidthOnly:
		return fmt.Sprintf("scale=%d:trunc(ow/a/2)*2", width)
	default:
		return ""
	}
}

// @TODO unused
func mkFn(max uint16) func(uint16) bool {
	return func(arg uint16) bool {
		return arg < max
	}
}

var even predicate = func(x uint16) bool { return x%2 == 0 }

var ErrNoArgs = errors.New("At least 1 arg must be provided")
var ErrTwoArgs = errors.New("2 args must be provided")
var ErrScaleOutOfRange = errors.New("ScaleType not recognized")

// A priori - arguments have been checked to not exceed the
//            respective dimensions of the input video
func (s ScaleType) Decode(args ...uint16) (string, error) {
	var c = len(args)

	if _, ok := scales[s]; !ok {
		return "", ErrScaleOutOfRange
	}

	switch {
	default:
		return "", ErrNoArgs

	case c == 1:
		if !even(args[0]) {
			return "", fmt.Errorf("%d is not even", args[0])
		}
		switch s {
		case WidthOnly:
			return s.interpolate(args[0], 0), nil
		case HeightOnly:
			return s.interpolate(0, args[0]), nil
		case WidthHeight:
			fallthrough
		default:
			return "", ErrTwoArgs
		}

	case c >= 2:
		switch {
		case !even(args[0]):
			return "", fmt.Errorf("%d is not even", args[0])
		case !even(args[1]):
			return "", fmt.Errorf("%d is not even", args[1])
		default:
			return s.interpolate(args[0], args[1]), nil
		}
	}
}
