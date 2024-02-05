package common

import "fmt"

const sshAuthorizedKeysPath = "%s/.ssh/authorized_keys"
const darwinHome = "/Users/ashwin"
const linuxHome = "/home/ashwin"

func SSHAuthorizedKeys(keys []string, os string) error {
	home := linuxHome
	if os == "darwin" {
		home = darwinHome
	}

	sshPath := fmt.Sprintf(sshAuthorizedKeysPath, home)
	return nil
}
