package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/ashwinath/anubis/pkg/logger"
)

const sshAuthorizedKeysPath = "%s/.ssh/authorized_keys"
const darwinHome = "/Users/ashwin"
const linuxHome = "/home/ashwin"
const newLine = "\n"
const sshdAnubisConfig = "/etc/ssh/sshd_config.d/anubis.conf"

func SSHAuthorizedKeys(keys []string, isDarwin bool) error {
	logger.Infof("Updating ssh authorized_keys")

	home := linuxHome
	if isDarwin {
		home = darwinHome
	}

	sshPath := fmt.Sprintf(sshAuthorizedKeysPath, home)

	// Create authorized_keys file if needed
	if _, err := os.Stat(sshPath); err != nil {
		if f, err := os.Create(sshPath); err != nil {
			return err
		} else {
			f.Close()
		}

		if err = os.Chown(sshPath, 1000, 1000); err != nil {
			return err
		}

		if err = os.Chmod(sshPath, 0644); err != nil {
			return err
		}
	}

	authKeyData, err := os.ReadFile(sshPath)
	if err != nil {
		return fmt.Errorf("could not read ssh authorized_key file at %s: %v", sshPath, err)
	}

	sshKeysSet := make(map[string]struct{})
	for _, key := range keys {
		found := false
		for _, line := range strings.Split(string(authKeyData), newLine) {
			if strings.Contains(line, key) {
				found = true
				break
			}
		}

		if !found {
			sshKeysSet[key] = struct{}{}
		}
	}

	var sb strings.Builder
	for key := range sshKeysSet {
		sb.WriteString(key)
		sb.WriteString(newLine)
	}

	f, err := os.OpenFile(sshPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("could not open %s for writing, %v", sshPath, err)
	}
	defer f.Close()

	keysToWrite := sb.String()

	if _, err := f.WriteString(keysToWrite); err != nil {
		return fmt.Errorf("error writing to %s: %v", sshPath, err)
	}

	if keysToWrite != "" {
		logger.Infof("wrote the following entries to %s, entries:\n%s", sshPath, keysToWrite)
	} else {
		logger.Infof("%s up to date, nothing to write", sshPath)
	}

	return nil
}

func HardenSSHDaemon() error {
	logger.Infof("Hardening sshd")

	if _, err := os.Stat(sshdAnubisConfig); err != nil {
		if f, err := os.Create(sshdAnubisConfig); err != nil {
			return err
		} else {
			f.Close()
		}

		if err := os.Chmod(sshdAnubisConfig, 0600); err != nil {
			return err
		}
	}

	configData, err := os.ReadFile(sshdAnubisConfig)
	if err != nil {
		return fmt.Errorf("could not read ssh authorized_key file at %s: %v", sshdAnubisConfig, err)
	}

	configsSet := make(map[string]struct{})
	for _, key := range []string{"PasswordAuthentication no", "PubkeyAuthentication yes"} {
		found := false
		for _, line := range strings.Split(string(configData), newLine) {
			if strings.Contains(line, key) {
				found = true
				break
			}
		}

		if !found {
			configsSet[key] = struct{}{}
		}
	}

	var sb strings.Builder
	for key := range configsSet {
		sb.WriteString(key)
		sb.WriteString(newLine)
	}

	f, err := os.OpenFile(sshdAnubisConfig, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("could not open %s for writing, %v", sshdAnubisConfig, err)
	}
	defer f.Close()

	keysToWrite := sb.String()

	if _, err := f.WriteString(keysToWrite); err != nil {
		return fmt.Errorf("error writing to %s: %v", sshdAnubisConfig, err)
	}

	if keysToWrite != "" {
		logger.Infof("wrote the following entries to %s, entries:\n%s", sshdAnubisConfig, keysToWrite)
	} else {
		logger.Infof("%s up to date, nothing to write", sshdAnubisConfig)
	}
	return nil
}
