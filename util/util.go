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
    "errors"
    "fmt"
    "log"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "time"
)

var (
    MissingProgramError = errors.New("program name is invalid")
    InvalidPath         = errors.New("bad path supplied")

    // Duration: 00:08:20
    regexStartTime = regexp.MustCompile(`^(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})$`)
)

// total - total time of the video
// ss - string representation of an instant to start the capture
func ParseStartTime(ss string, total time.Duration) (time.Duration, error) {
    if IsEmpty(ss) {
        log.Fatalf("%s not in format 00:00:00", ss)
    }
    //TODO
    time.ParseDuration("00h01m04s")
    return 0, nil
}

func ValidatePort(port int) error {
    if port < 1024 || port > 65535 {
        return fmt.Errorf("Port %d not in the range [1024, 65535]", port)
    }
    return nil
}

func ToPort(port int) string {
    return ":" + strconv.Itoa(port)
}

func expandTilde(path string) (string, error) {
    this_user, err := user.Current()
    if err != nil {
        return path, err
    }
    homeDirectory := this_user.HomeDir
    return strings.Replace(path, "~", homeDirectory, 1), nil
}

func SanitizeFile(path string) (string, error) {
    if IsEmpty(path) {
        return "", InvalidPath
    }

    candidateFile := filepath.Clean(path)
    candidateFile, err := expandTilde(candidateFile)
    if err != nil {
        return candidateFile, err
    }
    fi, err := os.Stat(candidateFile)
    if err != nil {
        return candidateFile, err
    }

    if fi.IsDir() {
        return candidateFile, InvalidPath
    }

    _, err = os.Open(candidateFile)
    if err != nil {
        return candidateFile, err
    }
    return candidateFile, nil
}

func IsEmpty(arg string) bool {
    return strings.TrimSpace(arg) == ""
}

func IsExistProgram(execName string) (bool, error) {
    if IsEmpty(execName) {
        return false, MissingProgramError
    }

    cmd := exec.Command(execName, "-h")
    if err := cmd.Start(); err != nil {
        return false, err
    }

    if err := cmd.Wait(); err != nil {
        return false, err
    }

    return true, nil
}
