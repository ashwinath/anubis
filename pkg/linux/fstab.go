package linux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ashwinath/anubis/pkg/logger"
)

const spacebar = " "
const newLine = "\n"
const comment = "#"
const fstabLocation = "/etc/fstab"

func FSTab(fstabEntries []string) error {
	// Create all directories first
	for _, entry := range fstabEntries {
		directory := strings.Split(entry, spacebar)[1]
		if _, err := os.Stat(directory); err != nil {
			if err := os.MkdirAll(directory, 0755); err != nil {
				return err
			}
		}
	}

	// write into fstab
	fsData, err := os.ReadFile(fstabLocation)
	if err != nil {
		return fmt.Errorf("could not read /etc/fstab: %v", err)
	}

	fstabSet := make(map[string]struct{})
	for _, entry := range fstabEntries {
		found := false
		for _, line := range strings.Split(string(fsData), newLine) {
			if strings.HasPrefix(line, comment) {
				continue
			}

			if strings.Contains(line, entry) {
				found = true
				break
			}

		}
		if !found {
			fstabSet[entry] = struct{}{}
		}
	}

	var sb strings.Builder
	for key := range fstabSet {
		sb.WriteString(key)
		sb.WriteString(newLine)
	}

	f, err := os.OpenFile(fstabLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("could not open %s for writing, %v", fstabLocation, err)
	}
	defer f.Close()

	entriesToWrite := sb.String()
	if _, err := f.WriteString(entriesToWrite); err != nil {
		return fmt.Errorf("error writing to %s: %v", fstabLocation, err)
	}

	if entriesToWrite != "" {
		logger.Infof("wrote the following entries to /etc/fstab, entries:\n%s", entriesToWrite)
	} else {
		logger.Infof("/etc/fstab up to date, nothing to write")
	}

	out, err := exec.Command("mount", "-a").CombinedOutput()
	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}

	return nil
}
