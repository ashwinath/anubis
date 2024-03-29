package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ashwinath/anubis/pkg/logger"
)

func Download(url, destination string, chownToUser bool) error {
	logger.Infof("downloading %s into %s", url, destination)
	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	// change owner
	if chownToUser {
		if err := os.Chown(destination, 1000, 1000); err != nil {
			return fmt.Errorf("could not chown %s, error: %s", destination, err)
		}
	}

	return nil
}
