package utils

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(source, destination string) error {
	inputFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)

	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	return nil
}
