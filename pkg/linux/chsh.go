package linux

import (
	"fmt"
	"os/exec"
)

func Chsh() error {
	if out, err := exec.Command("chsh", "--shell", "/usr/bin/zsh", "ashwin").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}
