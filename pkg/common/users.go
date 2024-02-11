package common

import (
	"fmt"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
)

func AddGroupToUser(groups []string) error {
	logger.Infof("adding user to groups: %v", groups)

	for _, group := range groups {
		out, err := exec.Command("usermod", "-aG", group, "ashwin").CombinedOutput()
		if err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	logger.Infof("done adding user to groups: %v", groups)

	return nil
}
