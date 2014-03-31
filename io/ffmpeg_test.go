package io

import (
    "github.com/javouhey/seneca/vendor/github.com/stretchr/testify/assert"
    "github.com/javouhey/seneca/util"
    "testing"
)

func TestVideoReader(t *testing.T) {
    vr := new(VideoReader)
    err := vr.reset2(4,
        func() string { return "/tmp" },
        func() string { return "/" },
        func() int64 { return int64(1234567) },
    )
    assert.Error(t, err, "")

    vr.Filename = "/home/putin/crimea.mp4"
    err = vr.reset2(4,
        func() string { return "/tmp" },
        func() string { return "/" },
        func() int64 { return int64(1234567) },
    )
    assert.Equal(t, vr.Gif, "crimea.gif")
    assert.Equal(t, vr.TmpDir, "/tmp/seneca/1234567")
    assert.Equal(t, vr.PngDir, "/tmp/seneca/1234567/p")
    assert.Equal(t, vr.TmpFile, "img-%04d.png")
}

var vfArgsFixtures = []struct {
    NeedScaling bool
    ScaleFilter string
    SpeedSpec   string
    out         string
    vf          bool
}{
    {true, "600:300", "", "600:300", true},
    {false, "", "", "", false},
    {true, "600:300", "1*PTS", "600:300,1*PTS", true},
    {false, "", "1*PTS", "1*PTS", true},
}

func TestVfArgs(t *testing.T) {
    fg := new(FrameGenerator)
    for i, tt := range vfArgsFixtures {
        a := util.NewArguments()
        a.SpeedSpec = tt.SpeedSpec
        a.NeedScaling = tt.NeedScaling
        a.ScaleFilter = tt.ScaleFilter
        b, s := fg.combineVf(a)
        if b != tt.vf || s != tt.out {
            t.Errorf("%d. Error out(%t), want %t // out(%q), want %q", i, b, tt.vf, s, tt.out)
        }
    }
}
