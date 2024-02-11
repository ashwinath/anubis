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
		log.Printf("error unmarshalling config file: %s\n", err)
	}

	switch *target {
	case "fedora":
		fedora(c)
	default:
	}
}

func fedora(c *config.Config) {
	logger.Infof("begin process for fedora server")

	err := linux.RegisterDNFRepository(c.Fedora.DNF.Repos)
	if err != nil {
		logger.Errorf("error registering fedora repositories, error: %s", err)
	}

	err = linux.InstallFedoraPackages(c.Fedora.DNF.Packages)
	if err != nil {
		logger.Errorf("error installing fedora packages, error: %s", err)
	}

	err = linux.StartAndEnableServices(c.Fedora.SystemdServices)
	if err != nil {
		logger.Errorf("error enabling systemd services, error: %s", err)
	}

	err = common.CloneDotfiles(c.Dotfiles.Repo.HTTP, c.Dotfiles.Repo.SSH, false)
	if err != nil {
		logger.Errorf("error cloning dotfiles, error: %s", err)
	}

	err = common.ConfigureDotFiles(c.Fedora.DotFiles, false)
	if err != nil {
		logger.Errorf("error cloning dotfiles, error: %s", err)
	}

	err = linux.Chsh()
	if err != nil {
		logger.Errorf("could not change shell, error: %s", err)
	}

	err = linux.InstallFedoraRpms(c.Fedora.DNF.RPM)
	if err != nil {
		logger.Errorf("error installing fedora rpms, error: %s", err)
	}

	err = linux.FSTab(c.Fedora.FSTab)
	if err != nil {
		logger.Errorf("error editing /etc/fstab, error: %s", err)
	}

	err = common.Pip(c.Fedora.Python.Packages)
	if err != nil {
		logger.Errorf("error running pip, error: %s", err)
	}

	err = common.SSHAuthorizedKeys(c.Fedora.SSHPublicKeys, false)
	if err != nil {
		logger.Errorf("error running pip, error: %s", err)
	}

	err = common.HardenSSHDaemon()
	if err != nil {
		logger.Errorf("error hardening sshd, error: %s", err)
	}

	err = common.InstallBinaries(c.Fedora.Binaries)
	if err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	err = common.DownloadAndRunBinaries(c.Fedora.RunBinaries)
	if err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	err = common.Sudoers()
	if err != nil {
		logger.Errorf("error writing sudoers, error: %s", err)
	}

	err = common.Golang(c.GoVersion, false)
	if err != nil {
		logger.Errorf("error installing golang, error: %s", err)
	}

	err = common.CompileYCM(false)
	if err != nil {
		logger.Errorf("error compiling YCM, error: %s", err)
	}

	err = common.Fzf(false)
	if err != nil {
		logger.Errorf("error installing fzf, error: %s", err)
	}

	err = common.Neovim(false)
	if err != nil {
		logger.Errorf("error installing neovim, error: %s", err)
	}
}
