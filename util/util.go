package util

import (
    "strings"
    "errors"
    "os/exec"
)

var (
    MissingProgramError = errors.New("program name is invalid")
)

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
