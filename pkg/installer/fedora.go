package installer

import (
	"sync"

	"github.com/ashwinath/anubis/pkg/common"
	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/linux"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

func Fedora(c *config.Config) {
	logger.Infof("begin processing for fedora server")

	// synchronous
	if err := linux.RegisterDNFRepository(c.Fedora.DNF.Repos); err != nil {
		logger.Errorf("error registering fedora repositories, error: %s", err)
	}

	if err := linux.EnableCoprPackages(c.Fedora.DNF.Copr); err != nil {
		logger.Errorf("error registering copr repositories, error: %s", err)
	}

	if err := linux.InstallFedoraPackages(c.Fedora.DNF.Packages); err != nil {
		logger.Errorf("error installing fedora packages, error: %s", err)
	}

	if err := linux.StartAndEnableServices(c.Fedora.SystemdServices); err != nil {
		logger.Errorf("error enabling systemd services, error: %s", err)
	}

	if err := common.AddGroupToUser(c.Fedora.GroupsToAddToUser); err != nil {
		logger.Errorf("error adding groups to user: %s", err)
	}

	if err := common.Golang(c.GoVersion, false); err != nil {
		logger.Errorf("error installing golang, error: %s", err)
	}

	if err := common.InstallBinaries(c.Fedora.Binaries); err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	if err := common.DownloadAndRunBinaries(c.Fedora.RunBinaries); err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	var wg sync.WaitGroup

	utils.Go(&wg, func() {
		if err := common.CloneDotfiles(c.Dotfiles.Repo.HTTP, c.Dotfiles.Repo.SSH, false); err != nil {
			logger.Errorf("error cloning dotfiles, error: %s", err)
		}

		if err := common.ConfigureDotFiles(c.Fedora.DotFiles, false); err != nil {
			logger.Errorf("error cloning dotfiles, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.Chsh(); err != nil {
			logger.Errorf("could not change shell, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.InstallFedoraRpms(c.Fedora.DNF.RPM); err != nil {
			logger.Errorf("error installing fedora rpms, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.FSTab(c.Fedora.FSTab); err != nil {
			logger.Errorf("error editing /etc/fstab, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Pip(c.Fedora.Python.Packages); err != nil {
			logger.Errorf("error running pip, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.SSHAuthorizedKeys(c.Fedora.SSHPublicKeys, false); err != nil {
			logger.Errorf("error running pip, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.HardenSSHDaemon(); err != nil {
			logger.Errorf("error hardening sshd, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Sudoers(); err != nil {
			logger.Errorf("error writing sudoers, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.CompileYCM(false); err != nil {
			logger.Errorf("error compiling YCM, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Fzf(false); err != nil {
			logger.Errorf("error installing fzf, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Neovim(false); err != nil {
			logger.Errorf("error installing neovim, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Alacritty(c.Fedora.AlacrittyTag, false); err != nil {
			logger.Errorf("error installing alacritty, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.UniversalCtags(false); err != nil {
			logger.Errorf("error installing universal ctags, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.NerdFonts(false); err != nil {
			logger.Errorf("error installing nerd fonts, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.Kubernetes(c.Fedora.Kubernetes); err != nil {
			logger.Errorf("error installing kubernetes, error: %s", err)
		}
	})

	wg.Wait()

	logger.Infof("done processing for fedora server")
}
