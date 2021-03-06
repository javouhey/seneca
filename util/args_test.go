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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var speedFixtures = []struct {
	in  string
	out string
	err bool
}{
	{"", "", true},
	{"placebo", "", false},
	{"veryfast", "setpts=1/3*PTS", false},
	{"faster", "setpts=1/2*PTS", false},
	{"slower", "setpts=2*PTS", false},
	{"veryslow", "setpts=3*PTS", false},
	{"ultrafast", "", true},
	{"blahblah", "", true},
}

func TestPreprocessSpeed(t *testing.T) {
	assert.Equal(t, "1", "1")
	for i, tt := range speedFixtures {
		a := NewArguments()
		err := preprocessSpeed(a, tt.in)
		switch {
		case err == nil:
			if tt.err {
				t.Errorf("%d. Error out(%t), want %t", i, err != nil, tt.err)
			}
		default:
			if !tt.err {
				t.Errorf("%d. Error out(%t), want %t", i, err != nil, tt.err)
			}
		}
		if a.SpeedSpec != tt.out {
			t.Errorf("%d. in(%q) => out(%q), want %q", i, tt.in, a.SpeedSpec, tt.out)
		}
	}
}

func TestPreprocessFrom(t *testing.T) {
	var zerot time.Time
	a := NewArguments()
	tz := TimeCode(zerot)
	assert.Equal(t, a.From, tz)
	assert.Equal(t, tz.String(), "00:00:00")

	assert.Error(t, preprocessFrom(a, "13:01:5"))

	assert.NoError(t, preprocessFrom(a, "13:01:05"))
	ti, _ := time.Parse("15:04:05", "13:01:05")
	assert.Equal(t, tz.String(), "00:00:00")
	assert.Equal(t, a.From, TimeCode(ti))
	assert.Equal(t, TimeCode(ti).String(), "13:01:05")
}

func TestPreprocessingScale(t *testing.T) {
	a := NewArguments()
	assert.NoError(t, preprocessScale(a, "_:600"))
	assert.Equal(t, a.NeedScaling, true)
	assert.Equal(t, a.ScaleFilter, "scale=trunc(oh*a/2)*2:600")

	a = NewArguments()
	assert.NoError(t, preprocessScale(a, "300:600"))
	assert.Equal(t, a.NeedScaling, true)
	assert.Equal(t, a.ScaleFilter, "scale=300:600")

	a = NewArguments()
	assert.NoError(t, preprocessScale(a, "300:_"))
	assert.Equal(t, a.NeedScaling, true)
	assert.Equal(t, a.ScaleFilter, "scale=300:trunc(ow/a/2)*2")
}

func TestScaleType(t *testing.T) {

	_, err := WidthOnly.Decode()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrNoArgs)
	}

	_, err = ScaleType(5).Decode(100)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrScaleOutOfRange)
	}

	_, err = ScaleType(5).Decode(100, 200)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrScaleOutOfRange)
	}

	_, err = ScaleType(5).Decode(500, 300, 150)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrScaleOutOfRange)
	}

	a, _ := WidthOnly.Decode(100)
	assert.Equal(t, a, "scale=100:trunc(ow/a/2)*2")

	_, err = WidthOnly.Decode(101)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), "101 is not even")
	}

	_, err = HeightOnly.Decode()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrNoArgs)
	}

	a, _ = HeightOnly.Decode(666)
	assert.Equal(t, a, "scale=trunc(oh*a/2)*2:666")

	_, err = HeightOnly.Decode(661)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), "661 is not even")
	}

	a, err = WidthHeight.Decode()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrNoArgs)
	}

	a, err = WidthHeight.Decode(640)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrTwoArgs)
	}

	a, _ = WidthHeight.Decode(640, 480)
	assert.Equal(t, a, "scale=640:480")

	_, err = WidthHeight.Decode(641, 480)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), "641 is not even")
	}

	_, err = WidthHeight.Decode(640, 481)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), "481 is not even")
	}

	a, _ = WidthHeight.Decode(1280, 760, 481, 211)
	assert.Equal(t, a, "scale=1280:760")
}
