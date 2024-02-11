package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const tmpAlacrittyDirLinux = "/home/ashwin/anubis/alacritty"
const tmpAlacrittyDirDarwin = "/Users/ashwin/anubis/alacritty"
const alacrittyURL = "https://github.com/alacritty/alacritty.git"

func Alacritty(version string, isDarwin bool) error {
	logger.Infof("installing alacritty")

	tmpAlacrittyDir := tmpAlacrittyDirLinux
	if isDarwin {
		tmpAlacrittyDir = tmpAlacrittyDirDarwin
	}

	if _, err := os.Stat(tmpAlacrittyDir); err != nil {
		if err := utils.GitClone(alacrittyURL, tmpAlacrittyDir); err != nil {
			return fmt.Errorf("could not clone alacritty")
		}
	}

	if err := utils.ExecAsUser(fmt.Sprintf("git -C %s fetch --tags", tmpAlacrittyDir)); err != nil {
		return fmt.Errorf("error fetching tags for alacritty, error: %s", err)
	}

	if err := utils.ExecAsUser(fmt.Sprintf("git -C %s checkout tags/%s", tmpAlacrittyDir, version)); err != nil {
		return fmt.Errorf("error checking out tag %s, error: %s", version, err)
	}

	if err := utils.ExecAsUser(fmt.Sprintf("cd %s; cargo build --release", tmpAlacrittyDir)); err != nil {
		return fmt.Errorf("could not compile alacritty, error: %s", err)
	}

	var err error
	if isDarwin {
		err = configureAlacrittyDarwin()
	} else {
		err = configureAlacrittyLinux()
	}

	if err != nil {
		logger.Infof("done installing alacritty")
	}

	return err
}

func configureAlacrittyLinux() error {
	// Terminfo
	cmd := exec.Command("tic", "-xe", "alacritty,alacritty-direct", "extra/alacritty.info")
	cmd.Dir = tmpAlacrittyDirLinux

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	if err := utils.CopyFile(fmt.Sprintf("%s/target/release/alacritty", tmpAlacrittyDirLinux), "/usr/local/bin/alacritty"); err != nil {
		return fmt.Errorf("failed to copy alacritty binary to /usr/local/bin, error: %s", err)
	}

	if err := utils.CopyFile(fmt.Sprintf("%s/extra/logo/alacritty-term.svg", tmpAlacrittyDirLinux), "/usr/share/pixmaps/Alacritty.svg"); err != nil {
		return fmt.Errorf("failed to copy alacritty icon to /usr/share/pixmaps, error: %s", err)
	}

	cmd = exec.Command("desktop-file-install", "extra/linux/Alacritty.desktop")
	cmd.Dir = tmpAlacrittyDirLinux
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	if out, err := exec.Command("update-desktop-database").CombinedOutput(); err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}

func configureAlacrittyDarwin() error {
	if err := utils.ExecAsUser(fmt.Sprintf("cd %s; make app", tmpAlacrittyDirDarwin)); err != nil {
		return fmt.Errorf("could not run make app, error: %s", err)
	}

	if err := utils.ExecAsUser(fmt.Sprintf("cd %s; cp -r target/release/osx/Alacritty.app /Applications/", tmpAlacrittyDirDarwin)); err != nil {
		return fmt.Errorf("could not copy app to /Applications/, error: %s", err)
	}

	return nil
}
