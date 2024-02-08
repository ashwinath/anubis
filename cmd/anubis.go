package main

import (
	"flag"
	"log"

	"github.com/ashwinath/anubis/pkg/common"
	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/linux"
	"github.com/ashwinath/anubis/pkg/logger"
)

func main() {
	target := flag.String("target", "", "target machine: fedora-server, fedora, darwin")
	configPath := flag.String("config", "", "config file location")
	flag.Parse()

	c, err := config.New(*configPath)
	if err != nil {
		log.Printf("error unmarshalling config file: %v\n", err)
	}

	switch *target {
	case "fedora":
		fedora(c)
	default:
	}
}

func fedora(c *config.Config) {
	logger.Infof("begin process for fedora server")

	err := linux.InstallFedoraPackages(c.Fedora.DNF.Packages)
	if err != nil {
		logger.Errorf("error installing fedora packages, error: %v", err)
	}

	err = common.CloneDotfiles(c.Dotfiles.Repo.HTTP, c.Dotfiles.Repo.SSH, false)
	if err != nil {
		logger.Errorf("error cloning dotfiles, error: %v", err)
	}

	err = common.ConfigureDotFiles(c.Fedora.DotFiles, false)
	if err != nil {
		logger.Errorf("error cloning dotfiles, error: %v", err)
	}

	err = linux.Chsh()
	if err != nil {
		logger.Errorf("could not change shell, error: %v", err)
	}

	err = linux.InstallFedoraRpms(c.Fedora.DNF.RPM)
	if err != nil {
		logger.Errorf("error installing fedora rpms, error: %v", err)
	}

	err = linux.FSTab(c.Fedora.FSTab)
	if err != nil {
		logger.Errorf("error editing /etc/fstab, error: %v", err)
	}

	err = common.Pip(c.Fedora.Python.Packages)
	if err != nil {
		logger.Errorf("error running pip, error: %v", err)
	}

	err = common.SSHAuthorizedKeys(c.Fedora.SSHPublicKeys, false)
	if err != nil {
		logger.Errorf("error running pip, error: %v", err)
	}

	err = common.HardenSSHDaemon()
	if err != nil {
		logger.Errorf("error hardening sshd, error: %v", err)
	}

	err = common.InstallBinaries(c.Fedora.Binaries)
	if err != nil {
		logger.Errorf("error installing binaries, error: %v", err)
	}

	err = common.DownloadAndRunBinaries(c.Fedora.RunBinaries)
	if err != nil {
		logger.Errorf("error installing binaries, error: %v", err)
	}
}
