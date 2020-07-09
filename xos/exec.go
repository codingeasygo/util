package xos

import (
	"os/exec"
	"runtime"
	"strings"
)

//RunCommand will exec command and get text from stdout
func RunCommand(command string) (text string, err error) {
	var bys []byte
	switch runtime.GOOS {
	case "windows":
		bys, err = exec.Command("cmd", "/C", command).Output()
	default:
		bys, err = exec.Command("bash", "-c", command).Output()
	}
	text = string(bys)
	return
}

//Run will exec command and get text from stdout
func Run(args ...string) (text string, err error) {
	return RunCommand(strings.Join(args, " "))
}
