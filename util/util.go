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
    "path/filepath"
    "strings"
    "errors"
    "os/exec"
    "os"
)

var (
    MissingProgramError = errors.New("program name is invalid")
    InvalidPath = errors.New("bad path supplied")
)

func SanitizeFile(path string) (string, error) {
    if IsEmpty(path) {
        return "", InvalidPath
    }

    candidateFile := filepath.Clean(path)
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
