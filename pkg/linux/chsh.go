package linux

import (
	"fmt"
	"os/exec"
)

func Chsh() error {
	out, err := exec.Command("chsh", "--shell", "/usr/bin/zsh", "ashwin").CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}
