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

	// enable ufw
	if out, err := exec.Command("systemctl", "enable", "--now", "ufw.service").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	// configure ufw
	if out, err := exec.Command("ufw", "default", "allow").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	if out, err := exec.Command("ufw", "enable", "-y").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	for _, port := range portsToAllow {
		if out, err := exec.Command("ufw", "allow", port).CombinedOutput(); err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	if out, err := exec.Command("ufw", "default", "deny").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	if out, err := exec.Command("ufw", "reload").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}
