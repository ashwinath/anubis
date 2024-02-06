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

func SSHAuthorizedKeys(keys []string, isDarwin bool) error {
	home := linuxHome
	if isDarwin {
		home = darwinHome
	}

	sshPath := fmt.Sprintf(sshAuthorizedKeysPath, home)

	// Create authorized_keys file if needed
	if _, err := os.Stat(sshPath); err != nil {
		_, err := os.Create(sshPath)
		if err != nil {
			return err
		}

		err = os.Chown(sshPath, 1000, 1000)
		if err != nil {
			return err
		}

		err = os.Chmod(sshPath, 0644)
		if err != nil {
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
