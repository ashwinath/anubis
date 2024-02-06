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
	case "fedora-server":
		fedoraServer(c)
	default:
	}
}

func fedoraServer(c *config.Config) {
	logger.Infof("begin process for fedora server")
	err := linux.InstallFedoraPackages(c.FedoraServer.DNF.Packages)
	if err != nil {
		logger.Errorf("error installing fedora packages, error: %v", err)
	}

	err = linux.InstallFedoraRpms(c.FedoraServer.DNF.RPM)
	if err != nil {
		logger.Errorf("error installing fedora rpms, error: %v", err)
	}

	err = linux.FSTab(c.FedoraServer.FSTab)
	if err != nil {
		logger.Errorf("error editing /etc/fstab, error: %v", err)
	}

	err = common.Pip(c.FedoraServer.Python.Packages)
	if err != nil {
		logger.Errorf("error running pip, error: %v", err)
	}
}
