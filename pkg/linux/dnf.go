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
	logger.Infof("installing packages %v", packages)
	commands := append([]string{"install", "-y"}, packages...)
	out, err := exec.Command("dnf", commands...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}

const tmpDir = "/tmp/anubis/rpm"

func InstallFedoraRpms(rpms []config.RPM) error {
	// Create temp folder
	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return err
		}
	}

	defer os.RemoveAll(tmpDir)

	// Download all files async first
	var wg sync.WaitGroup
	wg.Add(len(rpms))

	for _, rpm := range rpms {
		go func(u, d string) {
			defer wg.Done()
			utils.Download(u, d)
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

	return nil
}

func getRPMLocation(name string) string {
	return fmt.Sprintf("%s/%s", tmpDir, name)
}
