package linux

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

func InstallFedoraPackages(packages []string) error {
	logger.Infof("installing fedora packages %v", packages)

	commands := append([]string{"install", "-y"}, packages...)
	out, err := exec.Command("dnf", commands...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	logger.Infof("finished installing fedora packages")

	return nil
}

const tmpRPMDir = "/tmp/anubis/rpm"

func InstallFedoraRpms(rpms []config.RPM) error {
	logger.Infof("installing fedora rpms")

	// Create temp folder
	if _, err := os.Stat(tmpRPMDir); err != nil {
		if err := os.MkdirAll(tmpRPMDir, 0755); err != nil {
			return err
		}
	}

	defer os.RemoveAll(tmpRPMDir)

	// Download all files async first
	var wg sync.WaitGroup
	wg.Add(len(rpms))

	for _, rpm := range rpms {
		go func(u, d string) {
			defer wg.Done()
			utils.Download(u, d, false)
		}(rpm.URL, getRPMLocation(rpm.Name))
	}
	wg.Wait()

	// Install RPMs
	for _, rpm := range rpms {
		logger.Infof("installing rpm: %s", rpm.Name)
		out, err := exec.Command("dnf", "localinstall", "-y", getRPMLocation(rpm.Name)).CombinedOutput()
		if err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	logger.Infof("finished installing fedora rpms")

	return nil
}

func getRPMLocation(name string) string {
	return fmt.Sprintf("%s/%s", tmpRPMDir, name)
}

func RegisterDNFRepository(repos []string) error {
	logger.Infof("registering dnf repositories")

	for _, repo := range repos {
		out, err := exec.Command("dnf", "config-manager", "--add-repo", repo).CombinedOutput()
		if err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	logger.Infof("done registering dnf repositories")

	return nil
}
