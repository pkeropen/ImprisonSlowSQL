package pkg

import (
	"errors"
	"fmt"
	"os/exec"
)

func ExecCommand(command string) error {

	cmd := exec.Command("/bin/sh", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		errMsg := fmt.Sprintf("Execute Shell: %s failed with error: %s", command, err.Error())
		return errors.New(errMsg)
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	return nil
}
