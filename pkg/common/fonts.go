package common

import (
	"fmt"
	"os"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const fontLocationLinux = "/home/ashwin/nerd-fonts"
const fontLocationDarwin = "/Users/ashwin/nerd-fonts"
const nerdFontsRepoURL = "https://github.com/ryanoasis/nerd-fonts.git"

func NerdFonts(isDarwin bool) error {
	logger.Infof("Installing nerd fonts")

	loc := fontLocationLinux
	if isDarwin {
		loc = fontLocationDarwin
	}

	if _, err := os.Stat(loc); err != nil {
		if err := utils.GitClone(nerdFontsRepoURL, loc); err != nil {
			return fmt.Errorf("error cloning nerd fonts, error: %s", err)
		}
	}

	if err := utils.ExecAsUser(fmt.Sprintf("git -C %s pull --ff-only", loc)); err != nil {
		return fmt.Errorf("error pulling nerd fonts, error: %s", err)
	}

	if err := utils.ExecAsUser(fmt.Sprintf("cd %s; ./install.sh", loc)); err != nil {
		return fmt.Errorf("error installing nerd fonts, error: %s", err)
	}

	logger.Infof("Done installing nerd fonts")

	return nil
}
