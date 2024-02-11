package common

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const universalCtagsRepositoryURL = "https://github.com/universal-ctags/ctags.git"
const tmpUniversalCtagsDirLinux = "/home/ashwin/anubis/universalctags"
const tmpUniversalCtagsDirDarwin = "/Users/ashwin/anubis/universalctags"

func UniversalCtags(isDarwin bool) error {
	logger.Infof("installing universal ctags")

	tmpUniversalCtagsDir := tmpUniversalCtagsDirLinux
	if isDarwin {
		tmpUniversalCtagsDir = tmpUniversalCtagsDirDarwin
	}

	err := utils.GitClone(universalCtagsRepositoryURL, tmpUniversalCtagsDir)
	if err != nil {
		return fmt.Errorf("could not clone universal ctags")
	}

	defer os.RemoveAll(tmpUniversalCtagsDir)

	cmd := exec.Command("./autogen.sh")
	cmd.Dir = tmpUniversalCtagsDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	cmd = exec.Command("./configure")
	cmd.Dir = tmpUniversalCtagsDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	cmd = exec.Command("make")
	cmd.Dir = tmpUniversalCtagsDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	cmd = exec.Command("make", "install")
	cmd.Dir = tmpUniversalCtagsDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("done installing universal ctags")

	return nil
}
