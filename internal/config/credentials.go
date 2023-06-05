package config

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/pelletier/go-toml/v2"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Kubernetes struct {
	Credentials    TokenCredentials `toml:"credentials"`
	ServiceAddress string           `toml:"service_address"`
}
type Vault struct {
	Credentials    TokenCredentials `toml:"credentials"`
	ServiceAddress string           `toml:"service_address"`
}
type Harbor struct {
	Credentials    UsernamePasswordCredentials `toml:"credentials"`
	ServiceAddress string                      `toml:"service_address"`
}

type Credentials struct {
	Kubernetes Kubernetes `toml:"kubernetes"`
	Vault      Vault      `toml:"vault"`
	Harbor     Harbor     `toml:"harbor"`
}
type TokenCredentials struct {
	Token string `toml:"token"`
}
type UsernamePasswordCredentials struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func findFirstExistingFile(paths []string) (string, error) {
	for _, path := range paths {
		expandedPath, err := expandUser(path)
		if err != nil {
			return "", err
		}

		if _, err := os.Stat(expandedPath); err == nil {
			return path, nil
		}
	}
	return "", logs.Local().Errorf("No existing file found in paths: %v", paths)
}

func expandUser(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}
	return path, nil
}

func GetCredentials() (Credentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Credentials{}, err
	}

	c := Credentials{}
	credPaths := []string{
		"./automated.toml",
		"/etc/automated.toml",
		fmt.Sprintf("%s/.config/automated.toml", homeDir),
		fmt.Sprintf("%s/Projects/ChewedFeed/Automated/automated.toml", homeDir),
	}

	credPath, err := findFirstExistingFile(credPaths)
	if err != nil {
		return c, err
	}

	file, err := os.Open(credPath)
	if err != nil {
		return c, err
	}

	if err := toml.NewDecoder(file).Decode(&c); err != nil {
		return c, err
	}

	return c, nil
}
