package io_test

import (
	"testing"
)

// template test
func TestEmptyCheck(t *testing.T) {
	var a []int
	if nil != a {
		t.Errorf("should be nil for an uninitialized slice")
	}
}
