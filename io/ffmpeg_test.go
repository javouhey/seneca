package io

import (
    "github.com/javouhey/seneca/vendor/github.com/stretchr/testify/assert"
    "testing"
)

func TestVideoReader(t *testing.T) {
    vr := new(VideoReader)
    vr.reset2(4,
        func() string { return "/tmp" },
        func() string { return "/" },
        func() int64 { return int64(1234567) },
    )
    assert.Equal(t, vr.TmpDir, "/tmp/1234567", "")
    assert.Equal(t, vr.TmpFile, "img-%04d.png", "")
}
