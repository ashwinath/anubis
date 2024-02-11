package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const tmpAlacrittyDir = "/tmp/anubis/alacritty"
const alacrittyURL = "https://github.com/alacritty/alacritty.git"

func Alacritty(version string, isDarwin bool) error {
	logger.Infof("installing alacritty")

	if _, err := os.Stat(tmpAlacrittyDir); err != nil {
		err := utils.GitClone(alacrittyURL, tmpAlacrittyDir)
		if err != nil {
			return fmt.Errorf("could not clone alacritty")
		}
	}

	out, err := exec.Command(
		"git",
		"-C", tmpAlacrittyDir,
		"fetch", "--tags",
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	out, err = exec.Command(
		"git",
		"-C", tmpAlacrittyDir,
		"checkout", fmt.Sprintf("tags/%s", version),
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	err = utils.ExecAsUser(fmt.Sprintf("cd %s; cargo build --release", tmpAlacrittyDir))
	if err != nil {
		return fmt.Errorf("could not compile alacritty, error: %s", err)
	}

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
	cmd.Dir = tmpAlacrittyDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	err = utils.CopyFile(fmt.Sprintf("%s/target/release/alacritty", tmpAlacrittyDir), "/usr/local/bin/alacritty")
	if err != nil {
		return fmt.Errorf("failed to copy alacritty binary to /usr/local/bin, error: %s", err)
	}

	err = utils.CopyFile(fmt.Sprintf("%s/extra/logo/alacritty-term.svg", tmpAlacrittyDir), "/usr/share/pixmaps/Alacritty.svg")
	if err != nil {
		return fmt.Errorf("failed to copy alacritty icon to /usr/share/pixmaps, error: %s", err)
	}

	cmd = exec.Command("desktop-file-install", "extra/linux/Alacritty.desktop")
	cmd.Dir = tmpAlacrittyDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	out, err = exec.Command("update-desktop-database").CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}

func configureAlacrittyDarwin() error {
	err := utils.ExecAsUser(fmt.Sprintf("cd %s; make app", tmpAlacrittyDir))
	if err != nil {
		return fmt.Errorf("could not run make app, error: %s", err)
	}

	err = utils.ExecAsUser(fmt.Sprintf("cd %s; cp -r target/release/osx/Alacritty.app /Applications/", tmpAlacrittyDir))
	if err != nil {
		return fmt.Errorf("could not copy app to /Applications/, error: %s", err)
	}

	return nil
}
