package utils

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func ExecCommand(command string) (result string, e error) {

	cmd := exec.Command("/bin/sh", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		errMsg := fmt.Sprintf("Execute Shell: %s failed with error: %s", command, err.Error())
		return "", errors.New(errMsg)
	}
	log.Infof("Execute Shell:%s finished with output:\n%s", command, string(output))

	if output != nil {
		result = String(output)
	}

	return result, nil
}
