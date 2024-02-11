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
	if err := ExecAsUser(fmt.Sprintf("git clone %s %s", url, destination)); err != nil {
		return fmt.Errorf("could not clone to repository as user, error: %v", err)
	}

	return nil
}

func ExecAsUser(command string) error {
	logger.Infof("exec as user: %s", command)

	if out, err := exec.Command("su", "-", "ashwin", "-c", command).CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("done exec as user: %s", command)

	return nil
}
