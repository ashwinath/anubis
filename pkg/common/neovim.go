package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const neovimPlugURL = "https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim"
const vimPlugDestinationFolderLinux = "/home/ashwin/.local/share/nvim/site/autoload"
const vimPlugDestinationFolderDarwin = "/Users/ashwin/.local/share/nvim/site/autoload"

const ycmDirectoryLinux = "/home/ashwin/.vim/plugged/YouCompleteMe"
const ycmDirectoryDarwin = "/Users/ashwin/.vim/plugged/YouCompleteMe"

func Neovim(isDarwin bool) error {
	logger.Infof("installing neovim plugins")

	folder := vimPlugDestinationFolderLinux
	if isDarwin {
		folder = vimPlugDestinationFolderDarwin
	}

	if _, err := os.Stat(folder); err != nil {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	out, err := exec.Command("chown", "-R", "1000:1000", folder).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %s", string(out), err)
	}

	err = utils.Download(neovimPlugURL, fmt.Sprintf("%s/plug.vim", folder), true)
	if err != nil {
		return fmt.Errorf("could not download vim plug, error: %s", err)
	}

	err = utils.ExecAsUser("nvim +'PlugInstall --sync' +qa")
	if err != nil {
		return fmt.Errorf("could not run vim plug script, error %s", err)
	}

	logger.Infof("done installing neovim plugins")

	return nil
}

func CompileYCM(isDarwin bool) error {
	logger.Infof("Compiling YCM")

	ycmDirectory := ycmDirectoryLinux
	if isDarwin {
		ycmDirectory = ycmDirectoryDarwin
	}

	err := utils.ExecAsUser(
		fmt.Sprintf("cd %s; source ~/.zshrc; python install.py --go-completer --rust-completer", ycmDirectory),
	)
	if err != nil {
		return fmt.Errorf("Could not compile YCM, error: %s", err)
	}

	logger.Infof("Done compiling YCM")

	return nil
}
