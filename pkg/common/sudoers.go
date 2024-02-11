package common

import (
	"fmt"
	"os"

	"github.com/ashwinath/anubis/pkg/logger"
)

const sudoersConfigOverride = "Defaults secure_path = /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/var/lib/snapd/snap/bin:/usr/local/go/bin:/home/ashwin/go/bin\n"
const sudoersFileOverride = "/etc/sudoers.d/anubis"

func Sudoers() error {
	logger.Infof("writing sudoers override at %s", sudoersFileOverride)
	if _, err := os.Stat(sudoersFileOverride); err != nil {
		file, err := os.Create(sudoersFileOverride)
		if err != nil {
			return err
		}
		file.Close()
	}

	if err := os.WriteFile(sudoersFileOverride, []byte(sudoersConfigOverride), 0644); err != nil {
		return fmt.Errorf("Could not write to sudoers override file at %s", sudoersFileOverride)
	}

	logger.Infof("done writing sudoers override at %s", sudoersFileOverride)

	return nil
}
