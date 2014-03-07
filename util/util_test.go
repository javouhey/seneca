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

package util_test

import (
    "testing"
    "github.com/javouhey/seneca/util"
    "github.com/javouhey/seneca/vendor/github.com/stretchr/testify/assert"
)

func TestEmptyCheck(t *testing.T) {
    assert.True(t, util.IsEmpty("  "))
    assert.True(t, util.IsEmpty(""))
    assert.False(t, util.IsEmpty(" a "), "after trimming is length 1 is not empty")
    assert.False(t, util.IsEmpty("a"), "string of length 1 is not empty")
}
