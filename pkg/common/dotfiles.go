package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const dotConfigFolderLinux = "/home/ashwin/.config"
const dotConfigFolderDarwin = "/Users/ashwin/.config"
const homeFolderLinux = "/home/ashwin"
const homeFolderDarwin = "/Users/ashwin"
const zshrcLineLinux = "source ${HOME}/dotfiles/zsh/zshrc"
const zshrcLineDarwin = "source ${HOME}/dotfiles/zsh/mac.zshrc"
const darwinDotFilesLocation = "/Users/ashwin/dotfiles"
const linuxDotFilesLocation = "/home/ashwin/dotfiles"

func CloneDotfiles(gitURL string, gitSSHURL string, isDarwin bool) error {
	logger.Infof("Installing dotfiles")

	loc := linuxDotFilesLocation
	if isDarwin {
		loc = darwinDotFilesLocation
	}

	if _, err := os.Stat(loc); err == nil {
		if err := utils.ExecAsUser(fmt.Sprintf("git -C %s pull --ff-only", loc)); err != nil {
			return fmt.Errorf("error pulling dotfiles, error %s", err)
		}
	} else {
		if err := utils.GitClone(gitURL, loc); err != nil {
			return err
		}
	}

	logger.Infof("Done installing dotfiles")

	return nil
}

func ConfigureDotFiles(items []string, isDarwin bool) error {
	for _, item := range items {
		switch item {
		case "alacritty":
			applyConfigStyleFolder(item, "alacritty.toml", isDarwin)
		case "git":
			applyHomeDotfilesStyle(item, []string{"gitignore", "gitconfig"}, isDarwin)
		case "i3":
			applyConfigStyleFolder(item, "config", isDarwin)
		case "i3status-rust":
			applyConfigStyleFolder(item, "config.toml", isDarwin)
		case "nvim":
			applyConfigStyleFolder(item, "init.vim", isDarwin)
		case "ssh":
			ssh(isDarwin)
		case "tmux":
			applyHomeDotfilesStyle(item, []string{"theme.tmux", "tmux.conf"}, isDarwin)
		case "yabai":
			// TODO when I sort this out
		case "zsh":
			zshrc(isDarwin)
		}
	}
	return nil
}

func ssh(isDarwin bool) {
	logger.Infof("Install ssh config")

	home := homeFolderLinux
	if isDarwin {
		home = homeFolderDarwin
	}

	dotfilesLoc := linuxDotFilesLocation
	if isDarwin {
		dotfilesLoc = darwinDotFilesLocation
	}

	newLinkedLocation := fmt.Sprintf("%s/.ssh/config", home)
	if _, err := os.Stat(newLinkedLocation); err == nil {
		if err := os.Remove(newLinkedLocation); err != nil {
			logger.Errorf("error removing symlink at %s, error: %s", newLinkedLocation, err)
		}
	}

	if err := os.Symlink(fmt.Sprintf("%s/ssh/config", dotfilesLoc), newLinkedLocation); err != nil {
		logger.Errorf("error symlinking ssh config, error: %s", err)
	}

	if err := os.Chown(fmt.Sprintf("%s/ssh/config", dotfilesLoc), 1000, 1000); err != nil {
		logger.Errorf("error chowning ssh config, error: %s", err)
	}

	logger.Infof("Done install ssh config")
}

func zshrc(isDarwin bool) {
	logger.Infof("editing zshrc entry")

	home := homeFolderLinux
	if isDarwin {
		home = homeFolderDarwin
	}

	zshrcLoc := fmt.Sprintf("%s/.zshrc", home)

	zshrcData, err := os.ReadFile(zshrcLoc)
	if err != nil {
		logger.Errorf("could not read %s/.zshrc: %v", home, err)
		return
	}

	zshrcLine := zshrcLineLinux
	if isDarwin {
		zshrcLine = zshrcLineDarwin
	}

	found := false
	for _, line := range strings.Split(string(zshrcData), newLine) {
		if strings.Contains(line, zshrcLine) {
			found = true
			break
		}

	}
	if found {
		logger.Infof(".zshrc config is already updated, no op")
		return
	}

	f, err := os.OpenFile(zshrcLoc, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Errorf("could not open %s for writing, %v", zshrcLoc, err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(zshrcLine); err != nil {
		logger.Errorf("error writing to %s: %v", zshrcLoc, err)
		return
	}

	passwordsPath := fmt.Sprintf("%s/.passwords", home)
	if _, err := os.Stat(passwordsPath); err != nil {
		file, err := os.Create(passwordsPath)
		if err != nil {
			logger.Errorf("error creating file %s: %v", passwordsPath, err)
			return
		}
		file.Close()

		if err := os.Chown(passwordsPath, 1000, 1000); err != nil {
			logger.Errorf("error chowning file %s: %v", passwordsPath, err)
			return
		}
	}

	logger.Infof("Finished editing .zshrc")
}

func applyHomeDotfilesStyle(configFolder string, items []string, isDarwin bool) {
	logger.Infof("Installing %v", items)

	home := homeFolderLinux
	if isDarwin {
		home = homeFolderDarwin
	}

	dotfilesLoc := linuxDotFilesLocation
	if isDarwin {
		dotfilesLoc = darwinDotFilesLocation
	}

	for _, item := range items {
		newLinkedLocation := fmt.Sprintf("%s/.%s", home, item)
		if _, err := os.Stat(newLinkedLocation); err == nil {
			if err := os.Remove(newLinkedLocation); err != nil {
				logger.Errorf("error removing symlink at %s, error: %s", newLinkedLocation, err)
			}
		}

		if err := os.Symlink(fmt.Sprintf("%s/%s/%s", dotfilesLoc, configFolder, item), newLinkedLocation); err != nil {
			logger.Errorf("error symlinking %s config, error: %s", item, err)
			return
		}
		if err := os.Chown(fmt.Sprintf("%s/.%s", home, item), 1000, 1000); err != nil {
			logger.Errorf("error chowning %s config, error: %s", item, err)
			return
		}
	}

	logger.Infof("Done installing %v", items)
}

func applyConfigStyleFolder(configFolder string, item string, isDarwin bool) {
	logger.Infof("Installing %s", configFolder)

	loc := dotConfigFolderLinux
	if isDarwin {
		loc = dotConfigFolderDarwin
	}

	folder := fmt.Sprintf("%s/%s", loc, configFolder)
	if _, err := os.Stat(folder); err != nil {
		if err := os.MkdirAll(folder, 0755); err != nil {
			logger.Errorf("could not mkdirall %s: %v", folder, err)
			return
		}
	}

	dotfilesLoc := linuxDotFilesLocation
	if isDarwin {
		dotfilesLoc = darwinDotFilesLocation
	}

	newLinkedLocation := fmt.Sprintf("%s/%s/%s", loc, configFolder, item)
	if _, err := os.Stat(newLinkedLocation); err == nil {
		if err := os.Remove(newLinkedLocation); err != nil {
			logger.Errorf("error removing symlink at %s, error: %s", newLinkedLocation, err)
		}
	}

	if err := os.Symlink(fmt.Sprintf("%s/%s/%s", dotfilesLoc, configFolder, item), newLinkedLocation); err != nil {
		logger.Errorf("error symlinking %s config, error: %s", item, err)
	}

	// change owner
	if out, err := exec.Command("chown", "-R", "1000:1000", folder).CombinedOutput(); err != nil {
		logger.Errorf("error chowning %s, output: %s, error: %s", item, string(out), err)
		return
	}

	logger.Infof("Done installing alacritty config")
}
