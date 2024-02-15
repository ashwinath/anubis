package config

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Dotfiles              Dotfiles `yaml:"dotfiles"`
	Fedora                Fedora   `yaml:"fedora"`
	FedoraServerMaster    Fedora   `yaml:"fedoraServerMaster"`
	FedoraServerNonMaster Fedora   `yaml:"fedoraServerNonMasterServer"`
	GoVersion             string   `yaml:"goVersion"`
}

type Dotfiles struct {
	Repo GitRepo `yaml:"repo"`
}

type GitRepo struct {
	HTTP string `yaml:"http"`
	SSH  string `yaml:"ssh"`
}

type Fedora struct {
	AlacrittyTag         *string           `yaml:"alacrittyTag"`
	DNF                  DNF               `yaml:"dnf"`
	Python               Python            `yaml:"python"`
	Binaries             []Binary          `yaml:"binaries"`
	RunBinaries          []RunBinary       `yaml:"runBinaries"`
	Cargo                []string          `yaml:"cargo"`
	DotFiles             []string          `yaml:"dotfiles"`
	SSHPublicKeys        []string          `yaml:"sshPublicKeys"`
	FSTab                []string          `yaml:"fstab"`
	SystemdServices      []string          `yaml:"systemdServices"`
	GroupsToAddToUser    []string          `yaml:"groupsToAddToUser"`
	Kubernetes           *KubernetesConfig `yaml:"kubernetes"`
	FirewallPortsToAllow []string          `yaml:"string"`
}

type DNF struct {
	Copr     []string `yaml:"copr"`
	Repos    []string `yaml:"repos"`
	Packages []string `yaml:"packages"`
	RPM      []RPM    `yaml:"rpm"`
}

type RPM struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Python struct {
	Packages []string `yaml:"packages"`
}

type Binary struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Destination string `yaml:"destination"`
	Permissions uint32 `yaml:"permissions"`
	UID         *int   `yaml:"uid"`
	GID         *int   `yaml:"gid"`
}

type RunBinary struct {
	Name         string            `yaml:"name"`
	URL          string            `yaml:"url"`
	Flags        string            `yaml:"flags"`
	Env          map[string]string `yaml:"env"`
	AllowFailure bool              `yaml:"allowFailure"`
	ExecAsUser   bool              `yaml:"execAsUser"`
}

type KubernetesConfig struct {
	Version  string `yaml:"version"`
	IsMaster bool   `yaml:"isMaster"`
	MasterIP string `yaml:"masterIP"`
	NodeName string `yaml:"nodeName"`
}

func New(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	c := Config{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func NewFromGithubURL(url string) (*Config, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	c := Config{}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading config file returned non 200 status code, status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading res.Body from github config, error: %s", err)
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
