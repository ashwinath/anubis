package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ashwinath/anubis/pkg/logger"
)

const dotConfigFolderLinux = "/home/ashwin/.config"
const dotConfigFolderDarwin = "/Users/ashwin/.config"
const homeFolderLinux = "/home/ashwin"
const homeFolderDarwin = "/Users/ashwin"
const zshrcLine = "source ${HOME}/dotfiles/zsh/zshrc"
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

func ssh(isDarwin bool) error {
	logger.Infof("Install ssh config")

	home := homeFolderLinux
	if isDarwin {
		home = homeFolderDarwin
	}

	dotfilesLoc := linuxDotFilesLocation
	if isDarwin {
		dotfilesLoc = darwinDotFilesLocation
	}

	os.Symlink(
		fmt.Sprintf("%s/ssh/config", dotfilesLoc),
		fmt.Sprintf("%s/.ssh/config", home),
	)
	os.Chown(fmt.Sprintf("%s/ssh/config", dotfilesLoc), 1000, 1000)

	logger.Infof("Done install ssh config")
	return nil
}

func zshrc(isDarwin bool) error {
	logger.Infof("editing zshrc entry")

	home := homeFolderLinux
	if isDarwin {
		home = homeFolderDarwin
	}

	zshrcLoc := fmt.Sprintf("%s/.zshrc", home)

	zshrcData, err := os.ReadFile(zshrcLoc)
	if err != nil {
		return fmt.Errorf("could not read %s/.zshrc: %v", home, err)
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
		return nil
	}

	f, err := os.OpenFile(zshrcLoc, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("could not open %s for writing, %v", zshrcLoc, err)
	}
	defer f.Close()

	if _, err := f.WriteString(zshrcLine); err != nil {
		return fmt.Errorf("error writing to %s: %v", zshrcLoc, err)
	}

	logger.Infof("Finished editing .zshrc")

	return nil
}

func applyHomeDotfilesStyle(configFolder string, items []string, isDarwin bool) error {
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
		os.Symlink(
			fmt.Sprintf("%s/%s/%s", dotfilesLoc, configFolder, item),
			fmt.Sprintf("%s/.%s", home, item),
		)
		os.Chown(fmt.Sprintf("%s/.%s", home, item), 1000, 1000)
	}

	logger.Infof("Done installing %v", items)
	return nil
}

func applyConfigStyleFolder(configFolder string, item string, isDarwin bool) error {
	logger.Infof("Installing %s", configFolder)

	loc := dotConfigFolderLinux
	if isDarwin {
		loc = dotConfigFolderDarwin
	}

	folder := fmt.Sprintf("%s/%s", loc, configFolder)
	if _, err := os.Stat(folder); err != nil {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("could not mkdirall %s: %v", folder, err)
		}
	}

	dotfilesLoc := linuxDotFilesLocation
	if isDarwin {
		dotfilesLoc = darwinDotFilesLocation
	}

	os.Symlink(
		fmt.Sprintf("%s/%s/%s", dotfilesLoc, configFolder, item),
		fmt.Sprintf("%s/%s/%s", loc, configFolder, item),
	)

	// change owner
	out, err := exec.Command("chown", "-R", "1000:1000", folder).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("Done installing alacritty config")
	return nil
}
