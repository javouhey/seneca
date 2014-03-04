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

package io

import (
    "fmt"
    "bufio"
    "regexp"
    "strings"
)

const (
    strVideo = "Video:"
    strDuration = "Duration:"
)

var (
    duration = regexp.MustCompile(`Duration: (?P<duration>\d{2}:\d{2}:\d{2}.\d{2})`)
)

// Returns only the items that we need
func parse(str string) map[string]string{
    scanner := bufio.NewScanner(strings.NewReader(str))
    for scanner.Scan() {
        fmt.Println("=> ", scanner.Text())
    }
    return nil
}
