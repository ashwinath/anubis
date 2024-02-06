package common

import (
	"fmt"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
)

func CargoInstall(packages []string) error {
	logger.Infof("Installing cargo packages: %v", packages)

	commands := append([]string{"install"}, packages...)
	out, err := exec.Command("cargo", commands...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("Finished installing cargo packages")

	return nil
}
