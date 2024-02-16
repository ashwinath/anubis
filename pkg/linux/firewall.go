package linux

import (
	"fmt"
	"os/exec"
)

func ConfigureFirewall(portsToAllow []string) error {
	if len(portsToAllow) == 0 {
		// if it is not configured, best not touch anything
		return nil
	}

	// disable firewalld
	if out, err := exec.Command("systemctl", "stop", "firewalld").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	if out, err := exec.Command("systemctl", "disable", "firewalld").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}
