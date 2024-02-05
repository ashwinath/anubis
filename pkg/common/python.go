package common

import (
	"fmt"
	"log"
	"os/exec"
)

func Pip(pipPackages []string) error {
	pipBin := "pip3"
	if _, pip3Err := exec.LookPath("pip3"); pip3Err != nil {
		if _, pipErr := exec.LookPath("pip"); pipErr != nil {
			return fmt.Errorf("no pip installed")
		}
	}

	commands := append([]string{"install"}, pipPackages...)
	out, err := exec.Command(pipBin, commands...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}
	log.Printf("Installed pip packages.")

	return nil
}
