package io

import (
    "testing"
)

func TestEmptyCheck(t *testing.T) {
    var a []int
    if nil != a {
        t.Errorf("should be nil for an uninitialized slice")
    }
}
