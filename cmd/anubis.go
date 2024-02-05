package main

import (
	"flag"
	"log"

	"github.com/ashwinath/anubis/pkg/common"
	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/linux"
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
	err := linux.FSTab(c.FedoraServer.FSTab)
	if err != nil {
		log.Printf("error editing /etc/fstab: %v\n", err)
	}

	err = common.Pip(c.FedoraServer.Python.Packages)
	if err != nil {
		log.Printf("error running pip: %v", err)
	}
}
