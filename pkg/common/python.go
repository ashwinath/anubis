package common

import (
	"fmt"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
)

func Pip(pipPackages []string) error {
	logger.Infof("Installing pip packages: %v", pipPackages)

	pipBin := "pip3"
	if _, pip3Err := exec.LookPath("pip3"); pip3Err != nil {
		if _, pipErr := exec.LookPath("pip"); pipErr != nil {
			return fmt.Errorf("no pip installed")
		}
	}

	commands := append([]string{"install"}, pipPackages...)
	if out, err := exec.Command(pipBin, commands...).CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("Installed pip packages.")

	return nil
}
