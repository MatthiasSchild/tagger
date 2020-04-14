package main

import (
	"os/exec"
)

func execute(command string, args ...string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	return string(output), err
}
