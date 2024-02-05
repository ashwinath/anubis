package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("wont-crash", func(t *testing.T) {
		pwd, _ := os.Getwd()
		_, err := New(pwd + "/../../config.yaml")
		if err != nil {
			t.Errorf("error unmarshalling config: %v", err)
		}
	})
}
