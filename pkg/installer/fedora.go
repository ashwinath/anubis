package installer

import (
	"sync"

	"github.com/ashwinath/anubis/pkg/common"
	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/linux"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

func Fedora(c *config.Config, fedora config.Fedora, skipFonts bool) {
	logger.Infof("begin processing for fedora server")

	// synchronous
	if err := linux.RegisterDNFRepository(fedora.DNF.Repos); err != nil {
		logger.Errorf("error registering fedora repositories, error: %s", err)
	}

	if err := linux.EnableCoprPackages(fedora.DNF.Copr); err != nil {
		logger.Errorf("error registering copr repositories, error: %s", err)
	}

	if err := linux.InstallFedoraPackages(fedora.DNF.Packages); err != nil {
		logger.Errorf("error installing fedora packages, error: %s", err)
	}

	if err := linux.StartAndEnableServices(fedora.SystemdServices); err != nil {
		logger.Errorf("error enabling systemd services, error: %s", err)
	}

	if err := common.AddGroupToUser(fedora.GroupsToAddToUser); err != nil {
		logger.Errorf("error adding groups to user: %s", err)
	}

	if err := common.Golang(c.GoVersion, false); err != nil {
		logger.Errorf("error installing golang, error: %s", err)
	}

	if err := common.InstallBinaries(fedora.Binaries); err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	if err := common.DownloadAndRunBinaries(fedora.RunBinaries); err != nil {
		logger.Errorf("error installing binaries, error: %s", err)
	}

	var wg sync.WaitGroup

	utils.Go(&wg, func() {
		if err := common.CloneDotfiles(c.Dotfiles.Repo.HTTP, c.Dotfiles.Repo.SSH, false); err != nil {
			logger.Errorf("error cloning dotfiles, error: %s", err)
		}

		if err := common.ConfigureDotFiles(fedora.DotFiles, false); err != nil {
			logger.Errorf("error cloning dotfiles, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.Chsh(); err != nil {
			logger.Errorf("could not change shell, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.InstallFedoraRpms(fedora.DNF.RPM); err != nil {
			logger.Errorf("error installing fedora rpms, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := linux.FSTab(fedora.FSTab); err != nil {
			logger.Errorf("error editing /etc/fstab, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Pip(fedora.Python.Packages); err != nil {
			logger.Errorf("error running pip, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.SSHAuthorizedKeys(fedora.SSHPublicKeys, false); err != nil {
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
		if err := common.Fzf(false); err != nil {
			logger.Errorf("error installing fzf, error: %s", err)
		}
	})

	utils.Go(&wg, func() {
		if err := common.Neovim(false); err != nil {
			logger.Errorf("error installing neovim, error: %s", err)
		}

		if err := common.CompileYCM(false); err != nil {
			logger.Errorf("error compiling YCM, error: %s", err)
		}
	})

	if fedora.AlacrittyTag != nil {
		utils.Go(&wg, func() {
			if err := common.Alacritty(*fedora.AlacrittyTag, false); err != nil {
				logger.Errorf("error installing alacritty, error: %s", err)
			}
		})
	}

	utils.Go(&wg, func() {
		if err := common.UniversalCtags(false); err != nil {
			logger.Errorf("error installing universal ctags, error: %s", err)
		}
	})

	if !skipFonts {
		utils.Go(&wg, func() {
			if err := common.NerdFonts(false); err != nil {
				logger.Errorf("error installing nerd fonts, error: %s", err)
			}
		})
	}

	if fedora.Kubernetes != nil {
		utils.Go(&wg, func() {
			if err := linux.Kubernetes(*fedora.Kubernetes); err != nil {
				logger.Errorf("error installing kubernetes, error: %s", err)
			}
		})
	}

	utils.Go(&wg, func() {
		if err := linux.ConfigureFirewall(); err != nil {
			logger.Errorf("error configuring firewall, error: %s", err)
		}
	})

	wg.Wait()

	logger.Infof("done processing for fedora server")
}
