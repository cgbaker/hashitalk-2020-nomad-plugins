package python

import (
"bytes"
"fmt"
"os/exec"
"regexp"
"strings"
)

var pythonVersionCommand = []string{"python", "--version"}

func pythonVersionInfo() (version string, err error) {
	var out bytes.Buffer

	cmd := exec.Command(pythonVersionCommand[0], pythonVersionCommand[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to check python version: %v", err)
		return
	}

	version = parsePythonVersionOutput(out.String())
	return
}

func parsePythonVersionOutput(infoString string) (string) {
	infoString = strings.TrimSpace(infoString)

	lines := strings.Split(infoString, "\n")
	re := regexp.MustCompile("^Python (.+)$")
	if len(lines) == 0 {
		return ""
	}
	if sub := re.FindStringSubmatch(lines[0]); len(sub) == 2 {
		return sub[1]
	}
	return ""
}
