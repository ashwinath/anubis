package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
)

func GitClone(url, destination string) error {
	if _, err := os.Stat(destination); err == nil {
		// folder exists and do not need to clone
		return nil
	}

	// clone
	out, err := exec.Command("git", "clone", url, destination).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	// change owner
	out, err = exec.Command("chown", "-R", "1000:1000", destination).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}

func ExecAsUser(command string) error {
	logger.Infof("exec as user: %s", command)

	out, err := exec.Command("su", "-", "ashwin", "-c", command).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}
	logger.Infof("done exec as user: %s", command)

	return nil
}
