package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func ExecCmd(cmd string) []string {
	command := exec.Command("/bin/bash", "-c", cmd)
	okBuffer := bytes.NewBuffer(nil)
	errBuffer := bytes.NewBuffer(nil)
	command.Stdout = okBuffer
	command.Stderr = errBuffer
	var result string
	if err := command.Run(); err != nil {
		result = err.Error()
		result = fmt.Sprintf("error: %s %s", err.Error(), errBuffer.String())
	}
	result = fmt.Sprintf("%s%s", okBuffer.String(), result)
	reg := regexp.MustCompile("\n")
	resultSlice := reg.Split(result, -1)
	return resultSlice
}
