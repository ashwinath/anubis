package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
)

const darwinDotFilesLocation = "/Users/ashwin/dotfiles"
const linuxDotFilesLocation = "/home/ashwin/dotfiles"

func CloneDotfiles(gitURL string, gitSSHURL string, isDarwin bool) error {
	logger.Infof("Installing dotfiles")

	loc := linuxDotFilesLocation
	if isDarwin {
		loc = darwinDotFilesLocation
	}

	if _, err := os.Stat(loc); err == nil {
		// folder exists and do not need to clone
		return nil
	}

	// clone
	out, err := exec.Command("git", "clone", gitURL, loc).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	// Change remote
	out, err = exec.Command(
		"git",
		"-C", loc,
		"remote",
		"set-url", "origin",
		gitSSHURL,
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	// change owner
	out, err = exec.Command("chown", "-R", "1000:1000", loc).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("Done installing dotfiles")

	return nil
}
