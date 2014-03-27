package io

import (
    "github.com/javouhey/seneca/vendor/github.com/stretchr/testify/assert"
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
    assert.Equal(t, vr.TmpDir, "/tmp/1234567")
    assert.Equal(t, vr.TmpFile, "img-%04d.png")
}
