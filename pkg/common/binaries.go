package common

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const tmpBinariesDir = "/tmp/anubis/binaries"
const tmpRunBinariesDir = "/tmp/anubis/runbinaries"

func InstallBinaries(binaries []config.Binary) error {
	logger.Infof("Installing binaries")

	// Create temp folder
	if _, err := os.Stat(tmpBinariesDir); err != nil {
		if err := os.MkdirAll(tmpBinariesDir, 0755); err != nil {
			return err
		}
	}

	defer os.RemoveAll(tmpBinariesDir)

	// Download all files async first
	var wg sync.WaitGroup
	wg.Add(len(binaries))

	for _, bin := range binaries {
		go func(u, d string) {
			defer wg.Done()
			if err := utils.Download(u, d, false); err != nil {
				logger.Errorf("error downloading %s, error: %s", d, err)
			}
		}(bin.URL, getBinaryLocation(bin.Name))
	}
	wg.Wait()

	for _, bin := range binaries {
		if _, err := os.Stat(bin.Destination); err != nil {
			if err := os.MkdirAll(tmpBinariesDir, fs.FileMode(bin.Permissions)); err != nil {
				return err
			}
		}

		// copy downloaded file
		if err := utils.CopyFile(getBinaryLocation(bin.Name), bin.Destination); err != nil {
			return fmt.Errorf("error copying file from %s to %s, error: %s", getBinaryLocation(bin.Name), bin.Destination, err)
		}

		if bin.UID != nil && bin.GID != nil {
			// System folders should exist, so gid and uid should be present for these files
			// chown entire directory
			mainFolderSplit := strings.Split(bin.Destination, "/")[0:3]
			mainFolder := strings.Join(mainFolderSplit, "/")
			if out, err := exec.Command("chown", "-R", fmt.Sprintf("%d:%d", *bin.UID, *bin.GID), mainFolder).CombinedOutput(); err != nil {
				return fmt.Errorf("output: %s, error: %s", string(out), err)
			}
		}

		if err := os.Chmod(bin.Destination, fs.FileMode(bin.Permissions)); err != nil {
			return fmt.Errorf("could not chmod %s", bin.Destination)
		}
	}

	logger.Infof("Finished installing binaries")

	return nil
}

func getBinaryLocation(name string) string {
	return fmt.Sprintf("%s/%s", tmpBinariesDir, name)
}

func DownloadAndRunBinaries(binaries []config.RunBinary) error {
	logger.Infof("Downloading and running Binaries")
	// Create temp folder
	if _, err := os.Stat(tmpRunBinariesDir); err != nil {
		if err := os.MkdirAll(tmpRunBinariesDir, 0755); err != nil {
			return err
		}
	}

	defer os.RemoveAll(tmpRunBinariesDir)
	// Download all files async first
	var wg sync.WaitGroup
	wg.Add(len(binaries))

	for _, bin := range binaries {
		go func(u, d string) {
			defer wg.Done()
			if err := utils.Download(u, d, false); err != nil {
				logger.Errorf("error downloading %s, error: %s", d, err)
			}
		}(bin.URL, getRunBinaryLocation(bin.Name))
	}
	wg.Wait()

	for _, bin := range binaries {
		binaryLoc := getRunBinaryLocation(bin.Name)
		err := os.Chmod(binaryLoc, 0700)
		if err != nil {
			return fmt.Errorf("could not chmod binary %s: %v", bin.Name, err)
		}

		if bin.ExecAsUser {
			if err := utils.ExecAsUser(fmt.Sprintf("%s %s", binaryLoc, bin.Flags)); err != nil {
				return fmt.Errorf("could not exec as user, error: %s", err)
			}
		} else {
			cmd := exec.Command(binaryLoc, bin.Flags)
			cmd.Env = os.Environ()

			for key, value := range bin.Env {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
			}

			if out, err := cmd.CombinedOutput(); err != nil && !bin.AllowFailure {
				return fmt.Errorf("output: %s, error: %v", string(out), err)
			}
		}
	}

	logger.Infof("Finished downloading and running Binaries")

	return nil
}

func getRunBinaryLocation(name string) string {
	return fmt.Sprintf("%s/%s", tmpRunBinariesDir, name)
}
