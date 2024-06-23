package linux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

func InstallFedoraPackages(packages []string) error {
	logger.Infof("installing fedora packages %v", packages)

	if len(packages) == 0 {
		logger.Infof("no fedora packages to install")
		return nil
	}

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
			if err := utils.Download(u, d, false); err != nil {
				logger.Errorf("error downloading %s, error: %s", d, err)
			}
		}(rpm.URL, getRPMLocation(rpm.Name))
	}
	wg.Wait()

	// Install RPMs
	for _, rpm := range rpms {
		logger.Infof("installing rpm: %s", rpm.Name)
		if out, err := exec.Command("dnf", "localinstall", "-y", getRPMLocation(rpm.Name)).CombinedOutput(); err != nil {
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
		if out, err := exec.Command("dnf", "config-manager", "--add-repo", repo).CombinedOutput(); err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	logger.Infof("done registering dnf repositories")

	return nil
}

func EnableCoprPackages(repos []string) error {
	logger.Infof("registering copr repositories")

	for _, repo := range repos {
		if out, err := exec.Command("dnf", "copr", "enable", "-y", repo).CombinedOutput(); err != nil {
			return fmt.Errorf("output: %s, error: %v", string(out), err)
		}
	}

	logger.Infof("done registering copr repositories")

	return nil
}

func UpdateDNFRepo() error {
	if out, err := exec.Command("dnf", "upgrade", "--refresh", "-y").CombinedOutput(); err != nil {
		return fmt.Errorf("could not update dnf repositories, output: %s, error: %s", out, err)
	}
	return nil
}

func InstallDNFApp(app string, flags ...string) error {
	allFlags := []string{"install", "-y", app}
	allFlags = append(allFlags, flags...)
	if out, err := exec.Command("dnf", allFlags...).CombinedOutput(); err != nil {
		return fmt.Errorf("could install dnf app %s, output: %s, error: %s", app, out, err)
	}
	return nil
}

func GetLatestVersionDNFApp(app string, flags ...string) (string, error) {
	allFlags := []string{"list", "--showduplicates", app}
	allFlags = append(allFlags, flags...)
	out, err := exec.Command("dnf", allFlags...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could install list dnf app %s, output: %s, error: %s", app, out, err)
	}

	last := ""
	for _, s := range strings.Split(string(out), "\n") {
		if strings.Contains(s, "x86_64") {
			last = s
		}
	}
	return strings.Fields(last)[1], nil
}
