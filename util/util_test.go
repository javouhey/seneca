package util_test

import (
    "github.com/javouhey/seneca/util"
    "testing"
)

func TestEmptyCheck(t *testing.T) {
    if !util.IsEmpty("  ") {
        t.Errorf("2 spaces should be considered empty")
    }
    if !util.IsEmpty("") {
        t.Errorf("should be considered empty")
    }
    if util.IsEmpty(" a ") {
        t.Errorf("should NOT be considered empty")
    }
    if util.IsEmpty("a") {
        t.Errorf("should NOT be considered empty")
    }
}
