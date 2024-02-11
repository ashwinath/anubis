package common

import (
	"fmt"
	"os"

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
		return updateFzf(path)
	} else {
		return installFzf(path)
	}
}

func installFzf(path string) error {
	logger.Infof("Installing fzf")

	if err := utils.GitClone(fzfGitURL, path); err != nil {
		return err
	}

	if err := utils.ExecAsUser(fmt.Sprintf("%s/install --no-bash --no-fish --all", path)); err != nil {
		return err
	}

	logger.Infof("Done installing fzf")

	return nil
}

func updateFzf(path string) error {
	logger.Infof("Updating fzf")

	if err := utils.ExecAsUser(fmt.Sprintf("git -C %s pull --ff-only", path)); err != nil {
		return fmt.Errorf("error pulling from fzf, error: %s", err)
	}

	logger.Infof("Done updating fzf")

	return nil
}
