package common

import (
	"fmt"
	"strings"

	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

func CargoInstall(packages []string) error {
	logger.Infof("Installing cargo packages: %v", packages)

	if err := utils.ExecAsUser(fmt.Sprintf("cargo install %s", strings.Join(packages, " "))); err != nil {
		return fmt.Errorf("error running cargo install, error: %s", err)
	}

	logger.Infof("Finished installing cargo packages")

	return nil
}
