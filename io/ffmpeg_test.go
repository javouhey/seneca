package io

import (
    "github.com/javouhey/seneca/vendor/github.com/stretchr/testify/assert"
    "github.com/javouhey/seneca/util"
    "os"
    "path/filepath"
    "testing"
)

func TestVideoReader(t *testing.T) {
    var pathsep string = string(os.PathSeparator)

    vr := new(VideoReader)
    err := vr.reset2(4,
        func() string { return filepath.Join([]string{pathsep, "tmp"}...) },
        func() string { return string(os.PathSeparator) },
        func() int64 { return int64(1234567) },
    )
    assert.Error(t, err, "")

    src := []string{string(os.PathSeparator), "home", "putin", "crimea.mp4"}
    vr.Filename = filepath.Join(src...)
    err = vr.reset2(4,
        func() string { return filepath.Join([]string{pathsep, "tmp"}...) },
        func() string { return string(os.PathSeparator) },
        func() int64 { return int64(1234567) },
    )
    assert.Equal(t, vr.Gif, "crimea.gif")

    tmpdir := []string{string(os.PathSeparator), "tmp", "seneca", "1234567"}
    assert.Equal(t, vr.TmpDir, filepath.Join(tmpdir...))

    pngdir := []string{string(os.PathSeparator), "tmp", "seneca", "1234567", "p"}
    assert.Equal(t, vr.PngDir, filepath.Join(pngdir...))
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
