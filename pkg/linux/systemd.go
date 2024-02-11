package linux

import (
	"fmt"
	"os/exec"
)

func StartAndEnableServices(services []string) error {
	for _, svc := range services {
		out, err := exec.Command("systemctl", "enable", svc).CombinedOutput()
		if err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}

		out, err = exec.Command("systemctl", "start", svc).CombinedOutput()
		if err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	return nil
}
