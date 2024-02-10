package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const golangURLFormat = "https://go.dev/dl/go%s.%s.tar.gz"
const linuxArch = "linux-amd64"

const goTmpFolder = "/tmp/anubis/golang"

func Golang(goVersion string, isDarwin bool) error {
	if isDarwin {
		return installGolangDarwin(goVersion)
	}

	return installGolangLinux(goVersion)
}

func installGolangLinux(goVersion string) error {
	logger.Infof("Installing go binary on linux")

	if _, err := os.Stat(goTmpFolder); err != nil {
		if err := os.MkdirAll(goTmpFolder, 0755); err != nil {
			return err
		}
	}

	defer os.RemoveAll(goTmpFolder)

	tarFileDestination := fmt.Sprintf("%s/golang.tar.gz", goTmpFolder)
	err := utils.Download(
		fmt.Sprintf(golangURLFormat, goVersion, linuxArch),
		tarFileDestination,
		false,
	)

	if err != nil {
		return fmt.Errorf("could not download golang binary")
	}

	if _, err := os.Stat("/usr/local/go"); err == nil {
		err = os.Remove("/usr/local/go")
		if err != nil {
			return fmt.Errorf("could not remove old golang binary from /usr/local/go")
		}
	}

	// Extract file
	out, err := exec.Command(
		"tar", "-C", "/usr/local", "-xzf", tarFileDestination,
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf("output: %s, error: %s", string(out), err)
	}

	logger.Infof("Done installing go binary on linux")

	return nil
}

func installGolangDarwin(goVersion string) error {
	return nil
}