package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const fzfPathLinux = "/home/ashwin/.fzf"
const fzfPathDarwin = "/Users/ashwin/.fzf"
const fzfGitURL = "https://github.com/junegunn/fzf.git"

func Fzf(isDarwin bool) error {
	path := fzfPathLinux
	if isDarwin {
		path = fzfPathDarwin
	}

	if _, err := os.Stat(path); err == nil {
		updateFzf(path)
	} else {
		installFzf(path)
	}

	return nil
}

func installFzf(path string) error {
	logger.Infof("Installing fzf")

	err := utils.GitClone(fzfGitURL, path, true)
	if err != nil {
		return err
	}

	err = utils.ExecAsUser(fmt.Sprintf("%s/install --no-bash --no-fish --all", path))
	if err != nil {
		return err
	}

	logger.Infof("Done installing fzf")

	return nil
}

func updateFzf(path string) error {
	logger.Infof("Updating fzf")

	out, err := exec.Command(
		"git", "-C", path,
		"pull", "--ff-only",
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("Done updating fzf")

	return nil
}
