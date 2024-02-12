package main

import (
	"flag"
	"os"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/installer"
	"github.com/ashwinath/anubis/pkg/logger"
)

const tmpDir = "/tmp/anubis"
const anubisYAMLGithub = "https://raw.githubusercontent.com/ashwinath/dotfiles/master/anubis/anubis.yaml"

func main() {
	target := flag.String("target", "", "target machine: fedora-server, fedora, darwin")
	configPath := flag.String("config", "", "config file location")
	flag.Parse()

	var c *config.Config
	var err error
	if *configPath == "" {
		c, err = config.NewFromGithubURL(anubisYAMLGithub)
	} else {
		c, err = config.New(*configPath)
	}

	if err != nil {
		logger.Errorf("error unmarshalling/downloading config file: %s", err)
		return
	}

	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			logger.Errorf("error creating tmp directory, %s", err)
		}
	}

	switch *target {
	case "fedora":
		installer.Fedora(c, c.Fedora)
	case "fedora-server-master":
		installer.Fedora(c, c.FedoraServerMaster)
	case "fedora-server-non-master":
		installer.Fedora(c, c.FedoraServerNonMaster)
	default:
		logger.Infof("no op, no such target: %s", *target)
	}
}
