package main

import (
	"flag"
	"os"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/installer"
	"github.com/ashwinath/anubis/pkg/logger"
)

const tmpDir = "/tmp/anubis"

func main() {
	target := flag.String("target", "", "target machine: fedora-server, fedora, darwin")
	configPath := flag.String("config", "", "config file location")
	flag.Parse()

	c, err := config.New(*configPath)
	if err != nil {
		logger.Errorf("error unmarshalling config file: %s", err)
		return
	}

	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			logger.Errorf("error creating tmp directory, %s", err)
		}
	}

	switch *target {
	case "fedora":
		installer.Fedora(c)
	default:
	}
}
